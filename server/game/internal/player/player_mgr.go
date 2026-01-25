package player

import (
	"myplay/common/pb"

	"github.com/hechh/framework"
	"github.com/hechh/framework/actor"
	"github.com/hechh/framework/handler"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/uerror"
)

type PlayerMgr struct {
	actor.Actor
	mgr *actor.ActorMgr
}

func init() {
	handler.RegisterCmd((*PlayerMgr).Login)
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

func (d *PlayerMgr) Login(ctx framework.IContext, req *pb.LoginReq, rsp *pb.LoginRsp) error {
	if act := d.mgr.GetActor(ctx.GetActorId()); act != nil {
		return act.SendMsg(ctx.To("Player.Relogin"), req, rsp)
	}
	if req.Data == nil {
		return uerror.Err(pb.ErrorCode_ParameterInvalid, "参数错误")
	}
	usr := &Player{}
	usr.Init(ctx.GetId())
	if d.mgr.AddActor(usr) {
		return usr.SendMsg(ctx.To("Player.Login"), req, rsp)
	}
	return uerror.Err(pb.ErrorCode_ServiceHasStopped, "服务已经涨停")
}
