package playerfun

import (
	"myplay/common/pb"
	"myplay/server/game/internal/player/domain"

	"github.com/gogo/protobuf/proto"
)

type BagFun struct {
	*ParentFun
	PlayerFun
	data *pb.PlayerBaseData
}

func NewBagFun(parent *ParentFun) domain.IPlayerFun {
	return &BaseFun{ParentFun: parent}
}

func (d *BagFun) Load(msg *pb.PlayerData) error {
	d.data = msg.Base
	return nil
}

func (d *BagFun) Copy(data *pb.PlayerData) error {
	buf, err := proto.Marshal(d.data)
	if err != nil {
		return err
	}
	data.Base = &pb.PlayerBaseData{}
	return proto.Unmarshal(buf, data.Base)
}
