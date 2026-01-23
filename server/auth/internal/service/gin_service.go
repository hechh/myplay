package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"myplay/common/pb"
	"myplay/server/auth/internal/config"
	"myplay/server/auth/internal/middleware"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/hechh/library/async"
	"github.com/hechh/library/crypto"
	"github.com/hechh/library/uerror"
)

var (
	engine *gin.Engine
)

func Init() {
	engine = gin.Default()
	group := engine.Group("/auth")
	group.Use(middleware.Crypto(config.GetAesKey))
	{
		group.POST("/prelogin", middleware.Wrapper(Prelogin))
	}

	async.Go(func() {
		engine.Run(fmt.Sprintf("%s:%d", config.NodeCfg.Ip, config.NodeCfg.Port))
	})
}

func Mock() {
	gin.SetMode(gin.TestMode)
	engine = gin.Default()
	group := engine.Group("/auth")
	group.Use(middleware.Crypto(config.GetAesKey))
	{
		group.POST("/prelogin", middleware.Wrapper(Prelogin))
	}
}

func POST(router string, req *pb.AuthReq, rsp *pb.AuthRsp) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	body, err = crypto.AesEncrypto(body, []byte(config.GetAesKey()))
	if err != nil {
		return err
	}
	// 发送请求
	cli, _ := http.NewRequest("POST", router, bytes.NewBuffer(body))
	cli.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, cli)
	if w.Code != http.StatusOK {
		return uerror.Err(-1, "http请求错误：%d", w.Code)
	}

	// 解密包体
	if body, err = crypto.AesDecrypto(w.Body.Bytes(), []byte(config.GetAesKey())); err != nil {
		return err
	}
	return json.Unmarshal(body, rsp)
}
