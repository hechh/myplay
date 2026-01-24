package playerfun

import (
	"myplay/common/pb"
	"myplay/server/game/internal/player/domain"

	"github.com/gogo/protobuf/proto"
)

type BaseFun struct {
	*ParentFun
	PlayerFun
	data *pb.PlayerBaseData
}

func NewBaseFun(parent *ParentFun) domain.IPlayerFun {
	return &BaseFun{ParentFun: parent}
}

func (d *BaseFun) Load(msg *pb.PlayerData) error {
	d.data = msg.Base
	return nil
}

func (d *BaseFun) Copy(data *pb.PlayerData) error {
	buf, err := proto.Marshal(d.data)
	if err != nil {
		return err
	}
	data.Base = &pb.PlayerBaseData{}
	return proto.Unmarshal(buf, data.Base)
}
