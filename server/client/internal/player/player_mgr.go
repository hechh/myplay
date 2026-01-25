package player

import (
	"myplay/common/pb"
	"myplay/common/table/global_config"

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

func (d *PlayerMgr) Start() {
	d.mgr = new(actor.ActorMgr)
	d.mgr.Register(&Player{})
	d.mgr.Start()
	actor.Register(d.mgr)

	d.Actor.Register(d)
	d.Actor.Start()
	actor.Register(d)
}

func (d *PlayerMgr) Stop() {
	id := d.GetActorId()
	d.Actor.Stop()
	mlog.Infof("PlayerMgr.Close(%d)", id)
	d.mgr.Stop()
}

func (d *PlayerMgr) Close() {
	d.Stop()
	d.Done()
	d.Wait()
}

func (d *PlayerMgr) Remove(ctx framework.IContext) error {
	if usr := d.mgr.GetActor(ctx.GetActorId()); usr != nil {
		gc.Push(usr.(*Player).Close)
	}
	return nil
}

// 初始化玩家
func (d *PlayerMgr) Init(begin uint64, end uint64) error {
	cfg := global_config.SGet(0)
	if cfg == nil {
		return uerror.Err(pb.ErrorCode_TableNotFound, "配置不存在")
	}

	for i := begin; i <= end; i++ {
		uid := cfg.UidBeginValue + i
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
