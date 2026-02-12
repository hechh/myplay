package main

import (
	"flag"
	"myplay/common/pb"
	"myplay/message"
	"myplay/server/game/internal/config"
	"myplay/server/game/internal/player"

	"github.com/hechh/framework"
	"github.com/hechh/framework/actor"
	"github.com/hechh/framework/bus"
	"github.com/hechh/framework/cluster"
	"github.com/hechh/framework/context"
	"github.com/hechh/framework/gc"
	"github.com/hechh/framework/handler"
	"github.com/hechh/framework/packet"
	"github.com/hechh/framework/router"
	"github.com/hechh/library/async"
	"github.com/hechh/library/database"
	"github.com/hechh/library/fwatcher"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/myredis"
	"github.com/hechh/library/util"
)

func main() {
	var cfg string
	var nodeId int
	flag.IntVar(&framework.RunMode, "mode", 1, "启动模式")
	flag.StringVar(&cfg, "config", "config.yaml", "游戏配置文件")
	flag.IntVar(&nodeId, "id", 1, "服务ID")
	flag.Parse()

	// 加载配置
	util.Must(config.Load(cfg, int32(nodeId)))

	// 初始化日志库
	mlog.Init(framework.RunMode, config.NodeCfg.LogLevel, config.NodeCfg.LogPath, framework.GetSelfName())
	async.Except(mlog.Fatalf)

	mlog.Infof("初始化配置...")
	util.Must(fwatcher.Init(config.GameCfg.Common.TablePath))

	mlog.Infof("初始化数据库...")
	util.Must(database.Init(database.MysqlDriver, config.GameCfg.Mysql))

	mlog.Infof("初始化Redis...")
	util.Must(myredis.Init(config.GameCfg.Redis))

	mlog.Infof("初始化Snowflake...")
	//util.Must(snowflake.Init(framework.GetSelfType(), framework.GetSelfId()))

	mlog.Infof("初始化垃圾回收...")
	gc.Init()

	mlog.Infof("初始化路由...")
	router.Init(config.GameCfg.Router)

	mlog.Infof("初始化集群...")
	util.Must(cluster.Init(config.GameCfg.Cluster))

	mlog.Infof("初始化消息队列...")
	util.Must(bus.Init(config.GameCfg.Nats))
	util.Must(bus.SubscribeBroadcast(recv))
	util.Must(bus.SubscribeUnicast(recv))
	util.Must(bus.SubscribeReply(recv))

	mlog.Infof("注册Rpc...")
	message.Init()

	mlog.Infof("初始化PlayerMgr...")
	playerMgr := &player.PlayerMgr{}
	playerMgr.Init()

	mlog.Infof("服务启动成功...")
	util.Signal(func() {
		playerMgr.Close()
		bus.Close()
		cluster.Close()
		router.Close()
		gc.Close()
		mlog.Close()
	})
}

func recv(head *packet.Head, body []byte) {
	if head.ActorFunc == 0 {
		err := bus.Send(head, framework.Rpc(pb.NodeType_Gate, "Player.SendToClient", head.Id, body))
		if err != nil {
			mlog.Errorf("SendToClient失败: %v", err)
		}
		return
	}
	hh := handler.Get(head.ActorFunc)
	if hh == nil {
		mlog.Errorf("接口(%d)未注册", head.ActorFunc)
		return
	}
	if err := actor.Send(context.NewContext(head, hh.GetName()), body); err != nil {
		mlog.Errorf("Actor调用失败: %v", err)
	}
}
