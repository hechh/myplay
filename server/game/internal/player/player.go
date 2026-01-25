package player

import (
	"myplay/common/pb"
	"myplay/server/game/internal/player/domain"
	"myplay/server/game/internal/player/playerfun"
	"time"

	"github.com/hechh/framework"
	"github.com/hechh/framework/bus"
	"github.com/hechh/framework/context"
	"github.com/hechh/framework/handler"
	"github.com/hechh/framework/packet"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/uerror"
)

type Player struct {
	*playerfun.ParentFun
	heartTime int64
}

func init() {
	handler.RegisterCmd((*Player).Login)
	handler.RegisterCmd((*Player).Relogin)
	handler.RegisterCmd((*Player).Heart)
	handler.Register0(framework.EMPTY, (*Player).OnTick)
}

func (d *Player) Init(uid uint64) {
	d.ParentFun = playerfun.NewParentFun()
	d.Actor.Register(d)
	d.Actor.SetActorId(uid)
	d.Actor.Start()
}

func (d *Player) Close() {
	uid := d.GetActorId()
	d.Actor.Done()
	d.Actor.Wait()
	mlog.Infof("Player(%d)关闭成功", uid)
}

func (d *Player) OnTick(ctx framework.IContext) error {
	data := &pb.PlayerData{}
	flag := false
	d.ParentFun.Walk(func(fun domain.IPlayerFun) bool {
		if !fun.IsChange() {
			return true
		}
		if err := fun.Copy(data); err != nil {
			mlog.Errorf("数据拷贝失败: %v", err)
			return true
		}
		fun.Save()
		flag = true
		return true
	})
	if flag {
		bus.Send(&packet.Head{Id: d.GetActorId()}, framework.Rpc(pb.NodeType_Db, "PlayerMgr.Update", d.GetActorId(), data))
	}
	return nil
}

func (d *Player) Login(ctx framework.IContext, req *pb.LoginReq, rsp *pb.LoginRsp) error {
	// 加载模块
	for _, class := range playerfun.PlayerFunList {
		fun := playerfun.NewFunMap[class]
		d.ParentFun.RegisterFun(class, fun(d.ParentFun))
	}

	// 加载数据
	for _, class := range playerfun.PlayerFunList {
		fun := d.ParentFun.GetFun(class)
		if err := fun.Load(req.Data); err != nil {
			return err
		}
	}

	d.RegisterTimer(context.NewSimpleContext(ctx.GetActorId(), "Player.OnTick"), time.Second, -1)
	return bus.Send(ctx, framework.Rpc(pb.NodeType_Gate, "Player.LoginSuccess", ctx.GetId(), req))
}

func (d *Player) Relogin(ctx framework.IContext, req *pb.LoginReq, rsp *pb.LoginRsp) error {
	return bus.Send(ctx, framework.Rpc(pb.NodeType_Gate, "Player.LoginSuccess", ctx.GetId(), req))
}

// 心跳请求
func (p *Player) Heart(ctx framework.IContext, req *pb.HeartReq, rsp *pb.HeartRsp) error {
	now := time.Now().Unix()
	if p.heartTime <= 0 {
		p.heartTime = now
	}
	if now-p.heartTime >= framework.HeartTimeExpire {
		// todo: 剔除玩家
		return uerror.New(pb.ErrorCode_HeartTimeOver, "心跳超时")
	}

	p.heartTime = now
	rsp.BeginTime = req.BeginTime
	rsp.EndTime = now
	return nil
}
