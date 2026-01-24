package main

import (
	"flag"
	"myplay/message"
	"myplay/server/client/internal/config"
	"myplay/server/client/internal/player"

	"github.com/hechh/framework"
	"github.com/hechh/framework/gc"
	"github.com/hechh/library/async"
	"github.com/hechh/library/database"
	"github.com/hechh/library/fwatcher"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/util"
)

func main() {
	var cfg, mode string
	var nodeId int
	var begin, end int64
	flag.StringVar(&mode, "mode", "debug", "启动模式")
	flag.StringVar(&cfg, "config", "config.yaml", "游戏配置文件")
	flag.IntVar(&nodeId, "id", 1, "服务ID")
	flag.Int64Var(&begin, "begin", 1, "开始uid")
	flag.Int64Var(&end, "end", 1, "终止uid")
	flag.Parse()

	// 加载配置
	util.Must(config.Load(cfg, int32(nodeId)))

	// 初始化日志库
	mlog.Init(mode, config.NodeCfg.LogLevel, config.NodeCfg.LogPath, framework.GetSelfName())
	async.Except(mlog.Fatalf)

	mlog.Infof("初始化配置...")
	util.Must(fwatcher.Init(config.ClientCfg.Common.TablePath))

	mlog.Infof("初始化数据库...")
	util.Must(database.Init(database.MysqlDriver, config.ClientCfg.Mysql))

	// 注册rpc
	message.Init()

	// 初始化垃圾回收
	gc.Init()

	// 初始化玩家
	playerMgr := &player.PlayerMgr{}
	playerMgr.Start()
	util.Must(playerMgr.Init(uint64(begin), uint64(end)))

	util.Signal(func() {
		gc.Close()
		mlog.Close()
	})
}
