package main

import (
	"flag"
	"fmt"
	"myplay/common/pb"
	"myplay/message"
	"myplay/server/client/internal/config"
	"myplay/server/client/internal/player"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/hechh/framework"
	"github.com/hechh/framework/actor"
	"github.com/hechh/framework/context"
	"github.com/hechh/framework/gc"
	"github.com/hechh/framework/packet"
	"github.com/hechh/library/async"
	"github.com/hechh/library/database"
	"github.com/hechh/library/fwatcher"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/util"
	"github.com/spf13/cast"
)

var (
	playerMgr = &player.PlayerMgr{}
)

func main() {
	var cfg, mode string
	var nodeId int
	var uid int64
	flag.StringVar(&mode, "mode", "debug", "启动模式")
	flag.StringVar(&cfg, "config", "config.yaml", "游戏配置文件")
	flag.IntVar(&nodeId, "id", 1, "服务ID")
	flag.Int64Var(&uid, "uid", 1, "开始uid")
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
	playerMgr.Init()
	util.Must(playerMgr.Login(uint64(uid), uint32(nodeId)))

	r := gin.Default()
	r.POST("/send", send)
	r.POST("/login", login)
	async.Go(func() {
		r.Run(fmt.Sprintf("%s:%d", config.NodeCfg.Ip, config.NodeCfg.Port))
	})

	util.Signal(func() {
		playerMgr.Close()
		gc.Close()
		mlog.Close()
	})
}

func login(c *gin.Context) {
	uid := cast.ToUint64(c.PostForm("uid"))
	if uid <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": pb.ErrorCode_ParameterInvalid,
			"msg":  "uid参数错误",
		})
		return
	}
	nodeId := cast.ToUint32(c.PostForm("nodeId"))
	if nodeId <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": pb.ErrorCode_ParameterInvalid,
			"msg":  "gate节点参数错误",
		})
		return
	}
	if err := playerMgr.Login(uid, nodeId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": pb.ErrorCode_ParameterInvalid,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "登录成功")
}

func send(c *gin.Context) {
	// 解析 head
	headStr := c.PostForm("head")
	if len(headStr) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": pb.ErrorCode_ParameterInvalid,
			"msg":  "Head参数错误",
		})
		return
	}
	head := &packet.Head{}
	if err := proto.Unmarshal([]byte(headStr), head); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": pb.ErrorCode_ParameterInvalid,
			"msg":  err.Error(),
		})
		return
	}
	body := util.StringToBytes(c.PostForm("body"))
	if err := actor.Send(context.NewContext(head, "Player.Send"), body); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": pb.ErrorCode_ParameterInvalid,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "发送成功")
}
