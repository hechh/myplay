package player

import (
	"github.com/hechh/framework"
	"github.com/hechh/framework/actor"
	"github.com/hechh/framework/gc"
	"github.com/hechh/framework/handler"
	"github.com/hechh/library/mlog"
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
}

func (d *PlayerMgr) Remove(ctx framework.IContext) error {
	if usr := d.mgr.GetActor(ctx.GetActorId()); usr != nil {
		gc.Push(usr.(*Player).Close)
	}
	return nil
}

// 初始化玩家
func (d *PlayerMgr) Login(begin uint64, end uint64) error {
	for i := begin; i <= end; i++ {
		uid := i
		usr := &Player{}
		usr.Init(uid, 1)
		if err := usr.Login(); err != nil {
			return err
		}
		if !d.mgr.AddActor(usr) {
			mlog.Errorf("创建玩家(%d)失败", uid)
		}
	}
	return nil
}
