package player

import (
	"myplay/common/pb"
	"sync/atomic"
	"time"

	"github.com/hechh/framework"
	"github.com/hechh/framework/actor"
	"github.com/hechh/framework/bus"
	"github.com/hechh/framework/handler"
	"github.com/hechh/framework/packet"
	"github.com/hechh/framework/socket"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/uerror"
)

type Player struct {
	actor.Actor
	status    int32
	loginTime int64
	socketId  uint32
	extra     uint32
	version   uint32
}

func init() {
	handler.RegisterCmd((*Player).Login)                  // 登录
	handler.RegisterCmd((*Player).LoginSuccess)           // 登录成功
	handler.RegisterP1(framework.PROTO, (*Player).Kick)   // 剔除玩家
	handler.RegisterV1(framework.BYTES, (*Player).Handle) // 消息处理
}

func NewPlayer(head *packet.Head, now int64) *Player {
	return &Player{
		socketId:  head.SocketId,
		loginTime: now,
		extra:     head.Extra,
		version:   head.Version,
	}
}

func (d *Player) Init() {
	d.Actor.Register(d)
	d.Actor.Start()
}

func (d *Player) Close() {
	uid := d.GetActorId()
	d.Done()
	d.Wait()
	mlog.Infof("Player(%d)关闭成功", uid)
}

// 登录
func (d *Player) Login(ctx framework.IContext, req *pb.LoginReq, rsp *pb.LoginRsp) error {
	head := ctx.GetHead()
	if d.socketId > 0 && d.socketId != head.SocketId {
		// 关闭网络
		socket.Stop(d.socketId)
	}

	d.socketId = head.SocketId
	d.loginTime = time.Now().UnixMilli()
	d.extra = head.Extra
	d.version = head.Version

	// todo: 转发到db服务
	return nil
}

func (d *Player) LoginSuccess(ctx framework.IContext, req *pb.LoginReq, rsp *pb.LoginRsp) error {
	atomic.StoreInt32(&d.status, 1)
	ctx.AddDepth(1)
	// todo
	return nil
}

func (d *Player) Kick(ctx framework.IContext, event *pb.KickNotify) error {
	if event.Uid != d.GetActorId() || framework.GetSelfId() == event.NodeId {
		return nil
	}
	if d.loginTime > event.LoginTime {
		return nil
	}
	socket.Stop(d.socketId)
	return actor.SendMsgSimple(event.Uid, "PlayerMgr.Remove")
}

func (d *Player) Handle(ctx framework.IContext, body []byte) error {
	if atomic.LoadInt32(&d.status) <= 0 {
		socket.Stop(d.socketId)
		return uerror.Err(pb.ErrorCode_PlayerNotOnline, "玩家未登录成功")
	}

	head := ctx.GetHead()
	hh := handler.GetCmdRpc(head.Cmd)
	switch hh.GetNodeType() {
	case uint32(pb.NodeType_Gate):
		return d.Send(ctx.To(hh.GetName()), body)
	default:
		return bus.Send(&packet.Packet{Head: head, Body: body})
	}
}
