package test

import (
	"myplay/message"
	"myplay/server/db/internal/config"
	"path/filepath"
	"testing"

	"github.com/hechh/framework"
	"github.com/hechh/framework/bus"
	"github.com/hechh/framework/cluster"
	"github.com/hechh/framework/gc"
	"github.com/hechh/framework/router"
	"github.com/hechh/library/async"
	"github.com/hechh/library/database"
	"github.com/hechh/library/fwatcher"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/myredis"
	"github.com/hechh/library/util"
)

func TestMain(m *testing.M) {
	nodeId := 1
	cfg := "../../../configure/env/develop/config.yaml"

	// 加载配置
	util.Must(config.Load(cfg, int32(nodeId)))

	// 初始化日志库
	mlog.Init(1, config.NodeCfg.LogLevel, config.NodeCfg.LogPath, framework.GetSelfName())
	async.Except(mlog.Fatalf)

	mlog.Infof("初始化配置...")
	util.Must(fwatcher.Init("../../../configure/data"))

	mlog.Infof("初始化数据库...")
	util.Must(database.Init(database.MysqlDriver, config.DbCfg.Mysql))

	mlog.Infof("初始化Redis...")
	util.Must(myredis.Init(config.DbCfg.Redis))

	mlog.Infof("初始化Snowflake...")
	//util.Must(snowflake.Init(framework.GetSelfType(), framework.GetSelfId()))

	mlog.Infof("初始化垃圾回收...")
	gc.Init()

	mlog.Infof("初始化路由...")
	router.Init(config.DbCfg.Router)

	mlog.Infof("初始化集群...")
	util.Must(cluster.Init(config.DbCfg.Cluster))

	mlog.Infof("初始化消息队列...")
	util.Must(bus.Init(config.DbCfg.Nats))

	mlog.Infof("注册Rpc...")
	message.Init()

	m.Run()
}

func TestGlob(t *testing.T) {
	str := "output/data/aaa.conf"
	t.Log(filepath.Ext(str))

	files, err := filepath.Glob("../../../output/data/*.conf")
	t.Log(files, err)
}
