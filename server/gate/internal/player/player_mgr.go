package player

import (
	"myplay/common/dao/login_lock"
	"myplay/common/pb"
	"myplay/common/token"
	"myplay/server/gate/internal/config"
	"time"

	"github.com/hechh/framework"
	"github.com/hechh/framework/actor"
	"github.com/hechh/framework/bus"
	"github.com/hechh/framework/context"
	"github.com/hechh/framework/gc"
	"github.com/hechh/framework/handler"
	"github.com/hechh/framework/packet"
	"github.com/hechh/framework/router"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/uerror"
)

type PlayerMgr struct {
	actor.Actor
	mgr *actor.ActorMgr
}

func init() {
	handler.Register0(framework.EMPTY, (*PlayerMgr).Remove) // 删除玩家
	handler.RegisterP2(framework.PROTO, (*PlayerMgr).Login) // 登录
}

func (d *PlayerMgr) Init() {
	d.mgr = new(actor.ActorMgr)
	d.mgr.Register(&Player{})
	d.mgr.Start()
	actor.Register(d.mgr)

	d.Actor.Register(d)
	d.Actor.Start()
	actor.Register(d)
}

func (d *PlayerMgr) Close() {
	id := d.GetActorId()
	d.Done()
	d.Wait()
	mlog.Infof("PlayerMgr(%d)关闭成功", id)

	mgrId := d.mgr.GetActorId()
	d.mgr.Done()
	d.mgr.Wait()
	mlog.Infof("PlayerMgr.ActorMgr(%d)关闭成功", mgrId)
}

// 删除玩家
func (d *PlayerMgr) Remove(ctx framework.IContext) error {
	if usr := d.mgr.GetActor(ctx.GetActorId()); usr != nil {
		gc.Push(usr.(*Player).Close)
	}
	return nil
}

// 登录
func (d *PlayerMgr) Login(ctx framework.IContext, req *pb.LoginReq, rsp *pb.LoginRsp) error {
	// 解析 token
	ctx.Tracef("解析token")
	tok, err := token.ParseToken(req.Token, config.GateCfg.Common.TokenKey)
	if err != nil {
		return err
	}
	head := ctx.GetHead()
	head.Id = tok.Uid

	// 玩家已经在线
	ctx.Tracef("玩家是否已经在线")
	if usr := d.mgr.GetActor(tok.Uid); usr != nil {
		return usr.SendMsg(ctx.To("Player.Login"), req, rsp)
	}

	// 加载全局锁
	ctx.Tracef("全局加锁")
	if err := login_lock.Lock(head.Id, framework.GetSelfId(), 10*time.Second); err != nil {
		ctx.Tracef("全局加锁失败: %v", err)
		return err
	}

	// 无条件剔除其他节点登录
	now := time.Now()
	item := &pb.KickNotify{
		Uid:       tok.Uid,
		LoginTime: now.UnixMilli(),
		NodeId:    framework.GetSelfId(),
	}
	ctx.Tracef("推送剔除玩家广播：%v", item)
	err = bus.Broadcast(ctx.Copy(), framework.Rpc(pb.NodeType_Gate, "Player.Kick", tok.Uid, item))
	if err != nil {
		ctx.Tracef("推送剔除玩家广播失败：%v", err)
		return err
	}

	// 确保路由项已创建
	router.GetOrNew(uint32(pb.NodeType_Gate), tok.Uid)

	// 创建新玩家
	usr := &Player{}
	usr.Init(head, now.Unix())
	if d.mgr.AddActor(usr) {
		return usr.SendMsg(ctx.To("Player.Login"), req, rsp)
	}
	usr.Close()
	return uerror.New(pb.ErrorCode_ServiceHasStopped, "已经停止服务")
}

func (d *PlayerMgr) Handle(msg *packet.Packet) error {
	mlog.Tracef("接收到客户端消息：%v", msg)
	hh := handler.GetCmdRpc(msg.Head.Cmd)
	if hh == nil {
		return uerror.Err(pb.ErrorCode_CmdNotSupported, "Cmd(%d)接口未注册", msg.Head.Cmd)
	}
	switch msg.Head.Cmd {
	case uint32(pb.CMD_LOGIN_REQ):
		return d.Send(context.NewContext(msg.Head, "PlayerMgr.Login"), msg.Body)
	default:
		act := d.mgr.GetActor(msg.Head.Id)
		if act == nil {
			return uerror.Err(pb.ErrorCode_ActorIdNotExist, "玩家(%d)不存在", msg.Head.Id)
		}
		return act.Send(context.NewContext(msg.Head, "Player.Handle"), msg.Body)
	}
}
