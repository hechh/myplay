package player

import (
	"myplay/common/dao/account_data"
	"myplay/common/dao/player_data"
	"myplay/common/pb"

	"github.com/hechh/framework"
	"github.com/hechh/framework/actor"
	"github.com/hechh/framework/bus"
	"github.com/hechh/framework/handler"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/uerror"
)

type PlayerPool struct {
	actor.ActorPool
}

func init() {
	handler.RegisterCmd((*PlayerPool).Login)
	handler.Register0(framework.GOB, (*PlayerPool).Get)
	handler.RegisterP1(framework.PROTO, (*PlayerPool).Save)
}

func (d *PlayerPool) Init() {
	d.ActorPool.Register(d, 100)
	d.ActorPool.Start()
	actor.Register(d)
}

func (d *PlayerPool) Close() {
	id := d.GetActorId()
	d.Done()
	d.Wait()
	mlog.Infof("PlayerPool(%d)关闭成功", id)
}

// 加载数据
func (d *PlayerPool) Login(ctx framework.IContext, req *pb.LoginReq, rsp *pb.LoginRsp) error {
	data, err := player_data.Query(nil, ctx.GetId())
	if err != nil {
		return err
	}
	if data == nil {
		account, err := account_data.Query(nil, ctx.GetId())
		if err != nil {
			return err
		}
		if account == nil {
			return uerror.Err(pb.ErrorCode_MysqlFailed, "玩家(%d)账号不存在", ctx.GetId())
		}
		data = &pb.PlayerData{
			Uid: ctx.GetId(),
			Base: &pb.PlayerBaseData{
				Name:  account.Name,
				Email: account.Email,
				Phone: account.Phone,
			},
			Bag: &pb.PlayerBagData{},
		}
	}
	// 同步数据
	actor.SendMsgTo(ctx, "PlayerMgr.Load", data)
	req.Data = data
	return bus.Send(ctx, framework.Rpc(pb.NodeType_Game, "PlayerMgr.Login", ctx.GetId(), req))
}

func (d *PlayerPool) Get(ctx framework.IContext) error {
	data, err := player_data.Query(nil, ctx.GetId())
	if err != nil {
		return err
	}
	if data == nil {
		return actor.SendMsgTo(ctx, "PlayerMgr.Remove")
	}
	return actor.SendMsgTo(ctx, "PlayerMgr.Load", data)
}

func (d *PlayerPool) Save(ctx framework.IContext, data *pb.PlayerData) error {
	if err := player_data.Update(nil, data, "base", "bag"); err != nil {
		return err
	}
	return actor.SendMsgTo(ctx, "PlayerMgr.Load", data)
}
