package message

import (
	"myplay/common/pb"

	"github.com/hechh/framework"
	"github.com/hechh/framework/handler"
)

func Init() {
	handler.RegisterRpc2[pb.LoginReq, pb.LoginRsp](framework.PROTO, pb.NodeType_Gate, pb.CMD_LOGIN_REQ, "PlayerMgr.Login")     // 登录请求
	handler.RegisterRpc1[pb.KickNotify](framework.PROTO, pb.NodeType_Gate, pb.CMD_CMD_EMPTY, "Player.Kick")                    // 剔除玩家
	handler.RegisterRpc2[pb.LoginReq, pb.LoginRsp](framework.PROTO, pb.NodeType_Gate, pb.CMD_CMD_EMPTY, "Player.LoginSuccess") // 登录成功

	handler.RegisterRpc2[pb.LoginReq, pb.LoginRsp](framework.PROTO, pb.NodeType_Db, pb.CMD_CMD_EMPTY, "PlayerMgr.Login") // 登录
	handler.RegisterRpc1[pb.PlayerData](framework.PROTO, pb.NodeType_Db, pb.CMD_CMD_EMPTY, "PlayerMgr.Update")           // 更新PlayerData数据

	handler.RegisterRpc2[pb.LoginReq, pb.LoginRsp](framework.PROTO, pb.NodeType_Game, pb.CMD_CMD_EMPTY, "PlayerMgr.Login") // 登录
	handler.RegisterRpc2[pb.HeartReq, pb.HeartRsp](framework.PROTO, pb.NodeType_Game, pb.CMD_HEART_REQ, "Player.Heart")    // 心跳
}
