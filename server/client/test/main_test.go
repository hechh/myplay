package test

import (
	"myplay/common/pb"
	"myplay/server/client/internal/config"
	"myplay/server/client/internal/player"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hechh/framework"
	"github.com/hechh/framework/gc"
	"github.com/hechh/framework/handler"
	"github.com/hechh/library/async"
	"github.com/hechh/library/database"
	"github.com/hechh/library/fwatcher"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/util"
)

func TestClient(t *testing.T) {
	cfg := "../../../configure/env/develop/config.yaml"
	nodeId := 1
	begin := int64(1)
	end := int64(1)

	// 加载配置
	util.Must(config.Load(cfg, int32(nodeId)))

	// 初始化日志库
	mlog.Init("debug", config.NodeCfg.LogLevel, config.NodeCfg.LogPath, framework.GetSelfName())
	async.Except(mlog.Fatalf)

	// 初始化配置
	util.Must(fwatcher.Init("../../../configure/data"))

	mlog.Infof("初始化数据库...")
	util.Must(database.Init(database.MysqlDriver, config.ClientCfg.Mysql))

	// 初始化垃圾回收
	gc.Init()

	// 注册rpc
	//message.Init()

	// 初始化玩家
	playerMgr := &player.PlayerMgr{}
	playerMgr.Start()
	util.Must(playerMgr.Init(uint64(begin), uint64(end)))

	util.Signal(func() {
		gc.Close()
		mlog.Close()
	})
}

func TestMarshal(t *testing.T) {
	//message.Init()
	req := &pb.LoginReq{Token: "kaasdfasd;fkads"}
	rpc := handler.GetCmdRpc(uint32(pb.CMD_LOGIN_REQ))
	buf, err := rpc.Marshal(req)
	t.Log(err, buf)

	newReq := &pb.LoginReq{}
	err = rpc.Unmarshal(buf, newReq)
	t.Log(err, newReq)
}

func TestMarshal2(t *testing.T) {
	req := &pb.LoginReq{Token: "kaasdfasd;fkads"}
	buf, err := proto.Marshal(req)
	t.Log(err, buf)

	newReq := &pb.LoginReq{}
	err = proto.Unmarshal(buf, newReq)
	t.Log(err, newReq)
}
