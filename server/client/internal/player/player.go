package player

import (
	"fmt"
	"myplay/common/dao/account_data"
	"myplay/common/pb"
	"myplay/common/token"
	"myplay/server/client/internal/config"
	"myplay/server/client/internal/frame"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hechh/framework"
	"github.com/hechh/framework/actor"
	"github.com/hechh/framework/context"
	"github.com/hechh/framework/handler"
	"github.com/hechh/framework/packet"
	"github.com/hechh/framework/socket"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/uerror"
	"github.com/hechh/library/util"
	"github.com/spf13/cast"
)

type Player struct {
	actor.Actor
	mutex    sync.RWMutex
	client   *socket.SocketClient                // websocket连接
	sequence uint32                              // 发送序列号
	list     util.Map2[uint32, uint32, *Request] // 请求队列
	status   int32                               // 登录状态
	uid      uint64
	nodeId   uint32
}

func (d *Player) Init(uid uint64, nodeId uint32) {
	d.Actor.Register(d)
	d.Actor.SetActorId(uid)
	d.Actor.Start()
	d.uid = uid
	d.nodeId = nodeId
	d.list = make(util.Map2[uint32, uint32, *Request])
}

func (d *Player) Close() {
	uid := d.GetActorId()
	d.Done()
	d.Wait()
	mlog.Infof("Player(%d)关闭成功", uid)
}

func (d *Player) GetSequence() uint32 {
	return atomic.AddUint32(&d.sequence, 1)
}

func (d *Player) Login() error {
	// 生成玩家
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

	// 建立连接
	wsurl := fmt.Sprintf("ws://%s:%d/ws", config.NodeCfg.Ip, config.NodeCfg.Port)
	ws, _, err := websocket.DefaultDialer.Dial(wsurl, nil)
	if err != nil {
		return err
	}
	d.client = socket.NewSocketClient(socket.ConnWrapper(ws), &frame.Frame{}, 100*1024)

	// 发送登录请求
	sess := &pb.SessionData{
		Uid:       d.uid,
		Name:      usr.Name,
		LoginTime: time.Now().Unix(),
		Version:   "10.10.0",
		DeviceId:  cast.ToString(d.uid),
		Platform:  pb.Platform_Desktop,
	}
	str, err := token.GenToken(sess, config.ClientCfg.Common.TokenKey)
	if err != nil {
		return err
	}
	rpc := handler.GetCmdRpc(uint32(pb.CMD_LOGIN_REQ))
	if rpc == nil {
		return uerror.Err(-1, "接口(%s)未注册", pb.CMD_LOGIN_REQ)
	}
	return d.Write(rpc, &pb.LoginReq{Token: str})
}

func (d *Player) Write(rpc framework.IRpc, req any) error {
	buf, err := rpc.Marshal(req)
	if err != nil {
		return err
	}
	_, err = d.client.Write(&packet.Packet{
		Head: &packet.Head{
			DstNodeType: rpc.GetNodeType(),
			DstNodeId:   d.nodeId,
			Id:          d.uid,
			Cmd:         rpc.GetCmd(),
			Seq:         d.GetSequence(),
			ActorFunc:   rpc.GetCrc32(),
		},
		Body: buf,
	})
	return err
}

func (d *Player) Heart() {
	hh := handler.GetCmdRpc(uint32(pb.CMD_HEART_REQ))
	d.Write(hh, &pb.HeartReq{})
}

func (d *Player) loop() {
	for {
		pack, err := d.client.Read()
		if err != nil {
			mlog.Errorf("玩家(%d)连接断开: %v", d.uid, err)
			break
		}

		cmd := pack.Head.Cmd - (pack.Head.Cmd % 2)
		d.mutex.RLock()
		req, ok := d.list.Get(cmd, pack.Head.Seq)
		d.mutex.RUnlock()
		if !ok {
			continue
		}

		// 异步处理应答
		switch pack.Head.Cmd {
		case uint32(pb.CMD_LOGIN_REQ):
			// 是否登录成功
			if atomic.CompareAndSwapInt32(&d.status, 1, 1) {
				mlog.Errorf("玩家(%d)重复登录", d.uid)
				actor.SendMsg(context.NewSimpleContext(d.uid, "PlayerMgr.Remove"))
				return
			}
			atomic.StoreInt32(&d.status, 1)
			rsp, err := req.GetRsp(pack)
			if err != nil {
				mlog.Errorf("玩家(%d)登录失败: %v", d.uid, err)
				actor.SendMsg(context.NewSimpleContext(d.uid, "PlayerMgr.Remove"))
				return
			}
			d.RegisterTimer(context.NewSimpleContext(d.uid, "Player.Heart"), 3*time.Second, -1)
			mlog.Infof("玩家(%d)登录成功: %v", d.uid, rsp)
		default:
			rsp, err := req.GetRsp(pack)
			mlog.Infof("rsp:%v, error:%v", rsp, err)
		}

		// 删除请求记录
		d.mutex.Lock()
		d.list.Del(cmd, pack.Head.Seq)
		d.mutex.Unlock()
	}
}
