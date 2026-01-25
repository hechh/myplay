package player

import (
	"encoding/json"
	"myplay/common/pb"
	"time"

	"github.com/hechh/framework"
	"github.com/hechh/framework/actor"
	"github.com/hechh/framework/bus"
	"github.com/hechh/framework/handler"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/util"
)

type PlayerData struct {
	*pb.PlayerData
	status     int32
	change     int32
	updateTime int64
}

type PlayerMgr struct {
	actor.Actor
	datas map[uint64]*PlayerData
}

func init() {
	handler.Register0(framework.EMPTY, (*PlayerMgr).Remove)
	handler.Register0(framework.EMPTY, (*PlayerMgr).OnTick)
	handler.RegisterCmd((*PlayerMgr).Login)
	handler.RegisterP1(framework.PROTO, (*PlayerMgr).Load)
	handler.RegisterP1(framework.PROTO, (*PlayerMgr).Update)
}

func (d *PlayerMgr) Init() error {
	d.Actor.Register(d)
	d.Actor.Start()
	actor.Register(d)
	d.datas = make(map[uint64]*PlayerData)
	return d.RegisterTimer("PlayerMgr.OnTick", 5*time.Second, -1)
}

func (d *PlayerMgr) Close() {
	id := d.GetActorId()
	d.Actor.Done()
	for _, item := range d.datas {
		if item.status != 1 {
			buf, _ := json.Marshal(item.PlayerData)
			mlog.Errorf("增量数据保存失败: %s", string(buf))
			continue
		}
		actor.SendMsgSimple(item.Uid, "PlayerPool.Save", item.PlayerData)
	}
	d.Actor.Wait()
	mlog.Infof("PlayerMgr(%d)关闭成功", id)
}

func (d *PlayerMgr) Remove(ctx framework.IContext) error {
	delete(d.datas, ctx.GetId())
	return nil
}

func (d *PlayerMgr) OnTick(ctx framework.IContext) error {
	now := time.Now().Unix()
	for _, item := range d.datas {
		if item.status != 1 || now-item.updateTime < 5 {
			continue
		}
		item.updateTime = now
		actor.SendMsgSimple(item.Uid, "PlayerPool.Save", item.PlayerData)
	}
	return nil
}

// 登录
func (d *PlayerMgr) Login(ctx framework.IContext, req *pb.LoginReq, rsp *pb.LoginRsp) error {
	// 玩家数据存在
	item, ok := d.datas[ctx.GetId()]
	if !ok {
		item = &PlayerData{}
		d.datas[ctx.GetId()] = item
	}
	if item.status != 1 {
		return actor.SendMsgTo(ctx, "PlayerPool.Login", req, rsp)
	}
	req.Data = item.PlayerData
	return bus.Send(ctx, framework.Rpc(pb.NodeType_Game, "PlayerMgr.Login", ctx.GetId(), req))
}

// 加载玩家数据
func (d *PlayerMgr) Load(ctx framework.IContext, data *pb.PlayerData) error {
	item := d.datas[ctx.GetId()]
	if item.PlayerData == nil {
		item.PlayerData = data
	} else {
		item.Version = data.Version
		item.Base = util.Or(item.Base == nil, data.Base, item.Base)
		item.Bag = util.Or(item.Bag == nil, data.Bag, item.Bag)
	}
	item.status = 1
	item.change = 0
	item.updateTime = time.Now().Unix()
	return nil
}

// 更新数据
func (d *PlayerMgr) Update(ctx framework.IContext, data *pb.PlayerData) error {
	item, ok := d.datas[ctx.GetId()]
	if !ok {
		item = &PlayerData{PlayerData: data}
		d.datas[ctx.GetId()] = item
	}

	item.change++
	item.Base = util.Or(data.Base != nil, data.Base, item.Base)
	item.Bag = util.Or(data.Bag != nil, data.Bag, item.Bag)

	// 原始数据是否加载成功?
	if item.status != 1 {
		return actor.SendMsgTo(ctx, "PlayerPool.Get")
	}

	now := time.Now().Unix()
	if item.change >= 20 && now-item.updateTime >= 3 {
		item.updateTime = time.Now().Unix()
		return actor.SendMsgTo(ctx, "PlayerPool.Save", item.PlayerData)
	}
	return nil
}
