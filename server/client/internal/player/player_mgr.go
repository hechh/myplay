package player

import (
	"fmt"
	"myplay/common/dao/account_data"
	"myplay/common/pb"
	"myplay/server/client/internal/config"
	"time"

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

// 创建账号
func (d *PlayerMgr) get(uid uint64) (*pb.AccountData, error) {
	usr, err := account_data.Query(nil, uid)
	if err != nil {
		return nil, err
	}
	if usr != nil {
		return usr, nil
	}
	usr = &pb.AccountData{
		Uid:        uid,
		Name:       fmt.Sprintf("test%d", uid),
		Email:      fmt.Sprintf("%d@qq.com", uid),
		Phone:      fmt.Sprintf("135%d", uid),
		Password:   "12345",
		CreateTime: time.Now().Unix(),
		Platform:   pb.Platform_Desktop,
		LoginType:  pb.LoginType_Account,
	}
	if err := account_data.Insert(nil, usr); err != nil {
		return nil, err
	}
	return usr, nil
}

// 初始化玩家
func (d *PlayerMgr) Login(uid uint64, nodeId uint32) error {
	data, err := d.get(uid)
	if err != nil {
		return err
	}
	// 获取 ws 链接
	url, err := config.GetWsUrl(nodeId)
	if err != nil {
		return err
	}
	// 创建玩家客户端
	usr := &Player{}
	if err := usr.Init(uid, data.Name, url); err != nil {
		return err
	}
	if !d.mgr.AddActor(usr) {
		usr.Close()
		return uerror.Err(-1, "玩家(%d)登录失败", uid)
	}
	return nil
}
