package player

import (
	"fmt"
	"myplay/common/dao/account_data"
	"myplay/common/pb"
	"myplay/common/token"
	"myplay/server/client/internal/config"
	"myplay/server/client/internal/frame"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hechh/framework"
	"github.com/hechh/framework/actor"
	"github.com/hechh/framework/handler"
	"github.com/hechh/framework/packet"
	"github.com/hechh/framework/socket"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/uerror"
	"github.com/spf13/cast"
)

type Player struct {
	actor.Actor
	client   *socket.SocketClient // websocket连接
	sequence uint32               // 发送序列号
	status   int32                // 登录状态
	name     string
	uid      uint64
	nodeId   uint32
}

func init() {
	handler.Register0(framework.EMPTY, (*Player).OnTick)
	handler.RegisterV1(framework.BYTES, (*Player).Send)
}

func (d *Player) Init(uid uint64, nodeId uint32) {
	d.Actor.Register(d)
	d.Actor.SetActorId(uid)
	d.Actor.Start()
	d.uid = uid
	d.nodeId = nodeId
}

func (d *Player) Close() {
	uid := d.GetActorId()
	d.Done()
	d.Wait()
	mlog.Infof("Player(%d)关闭成功", uid)
}

func (d *Player) Login() error {
	// 生成玩家
	if err := d.build(); err != nil {
		return err
	}

	// 建立连接
	wsurl, err := config.GetWsUrl(d.nodeId)
	if err != nil {
		return err
	}
	ws, _, err := websocket.DefaultDialer.Dial(wsurl, nil)
	if err != nil {
		return err
	}
	d.client = socket.NewSocketClient(socket.ConnWrapper(ws), &frame.Frame{}, 100*1024)
	go d.loop()

	// 发送登录请求
	str, err := token.GenToken(&pb.SessionData{
		Uid:       d.uid,
		Name:      d.name,
		LoginTime: time.Now().Unix(),
		Version:   "10.10.0",
		DeviceId:  cast.ToString(d.uid),
		Platform:  pb.Platform_Desktop,
	}, config.ClientCfg.Common.TokenKey)
	if err != nil {
		return err
	}
	return d.write(uint32(pb.CMD_LOGIN_REQ), &pb.LoginReq{Token: str})
}

// 创建账号
func (d *Player) build() error {
	usr, err := account_data.Query(nil, d.uid)
	if err != nil {
		return err
	}
	if usr == nil {
		usr = &pb.AccountData{
			Uid:        d.uid,
			Name:       fmt.Sprintf("test%d", d.uid),
			Email:      fmt.Sprintf("%d@qq.com", d.uid),
			Phone:      fmt.Sprintf("135%d", d.uid),
			Password:   "12345",
			CreateTime: time.Now().Unix(),
			Platform:   pb.Platform_Desktop,
			LoginType:  pb.LoginType_Account,
		}
		if err := account_data.Insert(nil, usr); err != nil {
			return err
		}
	}
	d.name = usr.Name
	return nil
}

func (d *Player) write(cmd uint32, msg any) error {
	rpc := handler.GetCmdRpc(cmd)
	if rpc == nil {
		return uerror.Err(-1, "cmd(%d)未注册", cmd)
	}
	buf, err := rpc.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = d.client.Write(&packet.Packet{
		Head: &packet.Head{
			Id:        d.uid,
			Cmd:       cmd,
			Seq:       atomic.AddUint32(&d.sequence, 1),
			ActorFunc: rpc.GetCrc32(),
			ActorId:   d.uid,
		},
		Body: buf,
	})
	return err
}

func (d *Player) Send(ctx framework.IContext, body []byte) error {
	head := ctx.GetHead()
	head.Seq = atomic.AddUint32(&d.sequence, 1)
	_, err := d.client.Write(&packet.Packet{
		Head: head,
		Body: body,
	})
	return err
}

func (d *Player) OnTick(ctx framework.IContext) error {
	return d.write(uint32(pb.CMD_HEART_REQ), &pb.HeartReq{BeginTime: time.Now().Unix()})
}

func (d *Player) response(head *packet.Head, body []byte) (any, error) {
	hh := handler.GetCmdRpc(head.Cmd - 1)
	if hh == nil {
		return nil, uerror.Err(-1, "cmd(%d)未注册", head.Cmd)
	}
	rsp := hh.New(1)
	if err := hh.Unmarshal(body, rsp); err != nil {
		return nil, err
	}
	irsp, _ := rsp.(framework.IResponse)
	if msg := irsp.GetRspHead(); msg != nil {
		return nil, uerror.Err(msg.Code, msg.Msg)
	}
	return rsp, nil
}

func (d *Player) loop() {
	for {
		pack, err := d.client.Read()
		if err != nil {
			mlog.Errorf("玩家(%d)连接断开: %v", d.uid, err)
			return
		}

		// 等待登录成功返回
		if atomic.CompareAndSwapInt32(&d.status, 0, 0) {
			if pack.Head.Cmd != uint32(pb.CMD_LOGIN_RSP) {
				continue
			}
			rsp, err := d.response(pack.Head, pack.Body)
			if err != nil {
				mlog.Errorf("登录失败: %v", err)
				return
			}
			mlog.Infof("登录成功: %v", rsp)
			d.RegisterTimer("Player.OnTick", 3*time.Second, -1)
			atomic.AddInt32(&d.status, 1)
			continue
		}

		if pack.Head.Cmd%2 == 1 {
			rsp, err := d.response(pack.Head, pack.Body)
			if err != nil {
				mlog.Errorf("失败: %v", err)
			} else {
				mlog.Infof("成功：%v", rsp)
			}
		}
	}
}
