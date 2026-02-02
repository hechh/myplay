package player

import (
	"myplay/common/pb"
	"myplay/server/game/internal/player/domain"
	"myplay/server/game/internal/player/playerfun"
	"time"

	"github.com/hechh/framework"
	"github.com/hechh/framework/bus"
	"github.com/hechh/framework/handler"
	"github.com/hechh/framework/packet"
	"github.com/hechh/library/mlog"
)

type Player struct {
	*playerfun.ParentFun
	heartTime int64
}

func init() {
	handler.RegisterP2(framework.PROTO, (*Player).Login)
	handler.RegisterP2(framework.PROTO, (*Player).Relogin)
	handler.RegisterP2(framework.PROTO, (*Player).Heart)
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
	d.RegisterTimer("Player.OnTick", time.Second, -1)
	d.heartTime = time.Now().Unix()
	return bus.Send(ctx, framework.Rpc(pb.NodeType_Gate, "Player.LoginSuccess", ctx.GetId(), &pb.LoginReq{}))
}

func (d *Player) Relogin(ctx framework.IContext, req *pb.LoginReq, rsp *pb.LoginRsp) error {
	d.heartTime = time.Now().Unix()
	return bus.Send(ctx, framework.Rpc(pb.NodeType_Gate, "Player.LoginSuccess", ctx.GetId(), &pb.LoginReq{}))
}

// 心跳请求
func (d *Player) Heart(ctx framework.IContext, req *pb.HeartReq, rsp *pb.HeartRsp) error {
	d.heartTime = time.Now().Unix()
	rsp.BeginTime = req.BeginTime
	rsp.EndTime = d.heartTime
	return nil
}
