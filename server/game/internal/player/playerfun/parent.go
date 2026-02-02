package playerfun

import (
	"myplay/common/pb"
	"myplay/server/game/internal/player/domain"

	"github.com/hechh/framework/actor"
)

var (
	PlayerFunList = []pb.PlayerDataType{
		pb.PlayerDataType_BASE,
		pb.PlayerDataType_BAG,
	}
	NewFunMap = map[pb.PlayerDataType]func(*ParentFun) domain.IPlayerFun{
		pb.PlayerDataType_BASE: NewBaseFun,
		pb.PlayerDataType_BAG:  NewBagFun,
	}
)

type PlayerFun struct {
	ischange bool
}

func (d *PlayerFun) Change() {
	d.ischange = true
}

func (d *PlayerFun) Save() {
	d.ischange = false
}

func (d *PlayerFun) IsChange() bool {
	return d.ischange
}

type ParentFun struct {
	actor.Actor
	funs map[pb.PlayerDataType]domain.IPlayerFun
}

func NewParentFun() *ParentFun {
	return &ParentFun{funs: make(map[pb.PlayerDataType]domain.IPlayerFun)}
}

func (d *ParentFun) Walk(f func(domain.IPlayerFun) bool) {
	for _, data := range d.funs {
		if !f(data) {
			return
		}
	}
}

func (d *ParentFun) RegisterFun(tt pb.PlayerDataType, ff domain.IPlayerFun) {
	d.funs[tt] = ff
}

func (d *ParentFun) GetFun(tt pb.PlayerDataType) domain.IPlayerFun {
	return d.funs[tt]
}

func (d *ParentFun) GetBaseFun() *BaseFun {
	return d.GetFun(pb.PlayerDataType_BASE).(*BaseFun)
}

func (d *ParentFun) GetBagFun() *BagFun {
	return d.GetFun(pb.PlayerDataType_BAG).(*BagFun)
}
