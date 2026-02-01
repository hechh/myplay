package test

import (
	"myplay/common/pb"
	"myplay/server/auth/internal/config"
	"myplay/server/auth/internal/service"
	"testing"

	"github.com/hechh/framework"
	"github.com/hechh/library/async"
	"github.com/hechh/library/database"
	"github.com/hechh/library/fwatcher"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/myredis"
	"github.com/hechh/library/util"
)

func TestMain(m *testing.M) {
	// 加载配置
	util.Must(config.Load("../../../configure/env/develop/config.yaml", 1))

	// 初始化日志库
	mlog.Init(1, config.NodeCfg.LogLevel, config.NodeCfg.LogPath, framework.GetSelf().Name)
	async.Except(mlog.Fatalf)

	// 初始化配置
	util.Must(fwatcher.Init("../../../configure/data"))

	// 初始化数据库
	util.Must(database.Init(database.MysqlDriver, config.AuthCfg.Mysql))

	// 初始化redis
	util.Must(myredis.Init(config.AuthCfg.Redis))

	service.Mock()

	m.Run()
}

func TestPrelogin(t *testing.T) {
	req := &pb.AuthReq{
		Name:      "test01",
		Email:     "test@qq.com",
		Phone:     "123523452345",
		Platform:  pb.Platform_Ios,
		LoginType: pb.LoginType_Account,
		Ip:        "1.1.1.1",
		DeviceId:  "adfasdf012341lkjk",
		Version:   "10.0.0",
	}
	rsp := &pb.AuthRsp{}
	err := service.POST("/auth/prelogin", req, rsp)
	t.Log(err, rsp)
}

func Print(a string) any {
	return nil
}

func TestPritn(t *testing.T) {
	a := Print("asdfad")
	switch vv := a.(type) {
	case nil:
		t.Log("-0-----------nil")
	default:
		t.Log("----", vv)
	}
}
