package player

import (
	"github.com/hechh/framework"
	"github.com/hechh/framework/actor"
	"github.com/hechh/framework/gc"
	"github.com/hechh/framework/handler"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/uerror"
)

type PlayerMgr struct {
	actor.Actor
	mgr *actor.ActorMgr
}

func init() {
	handler.Register0(framework.EMPTY, (*PlayerMgr).Remove)
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
	mlog.Infof("PlayerMgr(%d)", id)

	mgrId := d.mgr.GetActorId()
	d.mgr.Done()
	d.mgr.Wait()
	mlog.Infof("PlayerMgr.ActorMgr(%d)关闭成功", mgrId)
}

func (d *PlayerMgr) Remove(ctx framework.IContext) error {
	if usr := d.mgr.GetActor(ctx.GetActorId()); usr != nil {
		gc.Push(usr.(*Player).Close)
	}
	return nil
}

// 初始化玩家
func (d *PlayerMgr) Login(uid uint64, nodeId uint32) error {
	usr := &Player{}
	usr.Init(uid, nodeId)
	// 登录
	if err := usr.Login(); err != nil {
		return err
	}
	if !d.mgr.AddActor(usr) {
		usr.Close()
		return uerror.Err(-1, "玩家(%d)登录失败", uid)
	}
	return nil
}
