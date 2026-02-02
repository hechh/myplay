package main

import (
	"flag"
	"myplay/server/auth/internal/config"
	"myplay/server/auth/internal/service"

	"github.com/hechh/framework"
	"github.com/hechh/library/async"
	"github.com/hechh/library/database"
	"github.com/hechh/library/fwatcher"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/myredis"
	"github.com/hechh/library/util"
)

func main() {
	var cfg, mode string
	var nodeId int
	flag.StringVar(&mode, "mode", "debug", "启动模式")
	flag.StringVar(&cfg, "config", "config.yaml", "游戏配置文件")
	flag.IntVar(&nodeId, "id", 1, "服务ID")
	flag.Parse()

	// 加载配置
	util.Must(config.Load(cfg, int32(nodeId)))

	// 初始化日志库
	mlog.Init(mode, config.NodeCfg.LogLevel, config.NodeCfg.LogPath, framework.GetSelfName())
	async.Except(mlog.Fatalf)

	// 初始化配置
	util.Must(fwatcher.Init(config.AuthCfg.Common.TablePath))

	// 初始化数据库
	util.Must(database.Init(database.MysqlDriver, config.AuthCfg.Mysql))

	// 初始化redis
	util.Must(myredis.Init(config.AuthCfg.Redis))

	// 注册路由
	service.Init()

	util.Signal(func() {
		myredis.Close()
		database.Close()
		mlog.Close()
	})
}
