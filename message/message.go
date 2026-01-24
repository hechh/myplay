package message

import (
	"myplay/common/pb"

	"github.com/hechh/framework/handler"
)

func Init() {
	handler.RegisterCmdRpc[pb.LoginReq, pb.LoginRsp](pb.NodeType_Gate, pb.CMD_LOGIN_REQ, "PlayerMgr.Login") // 登录请求
	handler.RegisterCmdRpc[pb.HeartReq, pb.HeartRsp](pb.NodeType_Gate, pb.CMD_HEART_REQ, "Playear.Heart")   // 心跳请求
}
