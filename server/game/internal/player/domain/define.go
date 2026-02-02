package domain

import "myplay/common/pb"

type IPlayerFun interface {
	Change()
	Save()
	IsChange() bool
	Load(*pb.PlayerData) error // 加载数据
	Copy(*pb.PlayerData) error // 保存数据
}
