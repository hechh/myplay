package test

import (
	"myplay/common/dao/router_data"
	"myplay/common/pb"
	"myplay/server/game/internal/config"
	"sync"
	"testing"

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
	"github.com/hechh/library/uerror"
	"github.com/hechh/library/util"
)

func TestGate(t *testing.T) {
	cfg := "../../../configure/env/develop/config.yaml"
	nodeId := 1

	// 加载配置
	util.Must(config.Load(cfg, int32(nodeId)))

	// 初始化日志库
	mlog.Init("debug", config.NodeCfg.LogLevel, config.NodeCfg.LogPath, framework.GetSelfName())
	async.Except(mlog.Fatalf)

	mlog.Infof("初始化配置...")
	util.Must(fwatcher.Init("../../../configure/data"))

	mlog.Infof("初始化数据库...")
	util.Must(database.Init(database.MysqlDriver, config.GameCfg.Mysql))

	mlog.Infof("初始化Redis...")
	util.Must(myredis.Init(config.GameCfg.Redis))

	mlog.Infof("初始化Snowflake...")
	//util.Must(snowflake.Init(framework.GetSelfType(), framework.GetSelfId()))

	mlog.Infof("初始化垃圾回收...")
	gc.Init()

	mlog.Infof("初始化路由...")
	router.Init(config.NodeCfg, router_data.SaveRouter)

	mlog.Infof("初始化集群...")
	util.Must(cluster.Init(config.GameCfg.Etcd))

	mlog.Infof("初始化消息队列...")
	util.Must(bus.Init(config.GameCfg.Nats))

	mlog.Infof("注册Rpc...")
	//message.Init()

	mlog.Infof("服务启动成功...")
	util.Signal(func() {
		bus.Close()
		cluster.Close()
		router.Close()
		gc.Close()
		mlog.Close()
	})
}

// 处理客户端消息
func handle(msg *packet.Packet) error {
	h := handler.GetRpc(msg.Head.DstNodeType, msg.Head.ActorFunc)
	if h == nil {
		return uerror.Err(pb.ErrorCode_RpcNotSupported, "接口没有注册:%v", msg.Head)
	}
	switch h.GetNodeType() {
	case uint32(pb.NodeType_Gate):
		return actor.Send(context.NewContext(msg.Head, h.GetName()), msg.Body)
	default:
		return bus.Send(msg)
	}
}

func TestMgr(t *testing.T) {
	aa := make(chan int32, 10)
	aa <- 123
	aa <- 12
	aa <- 1
	b := sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		b.Add(1)
		go func() {
			defer t.Log("=====finish=====")
			for item := range aa {
				t.Log(item)
			}
			b.Done()
		}()
	}
	close(aa)
	b.Wait()
	t.Log("----------")
	//t.Log(<-aa, <-aa, <-aa, <-aa, <-aa)
}
