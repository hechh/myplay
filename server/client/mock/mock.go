package mock

import (
	"fmt"
	"myplay/server/client/internal/config"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hechh/framework/handler"
	"github.com/hechh/framework/packet"
	"github.com/hechh/library/uerror"
	"github.com/spf13/cast"
)

func Init(cfg string, nodeId uint32) error {
	return config.Load(cfg, int32(nodeId))
}

func Login(uid uint64, nodeId uint32) error {
	formData := url.Values{}
	formData.Add("uid", cast.ToString(uid))
	formData.Add("nodeId", cast.ToString(nodeId))

	// 发送请求
	targetUrl := fmt.Sprintf("http://%s:%d/login", config.NodeCfg.Ip, config.NodeCfg.Port)
	cli, err := http.NewRequest("POST", targetUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	cli.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(cli)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func Send(uid uint64, aid uint64, cmd uint32, msgs ...any) error {
	rpc := handler.GetCmdRpc(cmd)
	if rpc == nil {
		return uerror.Err(-1, "cmd(%d)未注册", cmd)
	}
	buf, err := rpc.Marshal(msgs...)
	if err != nil {
		return err
	}
	head := &packet.Head{
		Id:        uid,
		Cmd:       cmd,
		ActorFunc: rpc.GetCrc32(),
		ActorId:   aid,
	}
	headBuf, err := proto.Marshal(head)
	if err != nil {
		return err
	}

	formData := url.Values{}
	formData.Add("head", string(headBuf))
	formData.Add("body", string(buf))

	// 发送请求
	targetUrl := fmt.Sprintf("http://%s:%d/send", config.NodeCfg.Ip, config.NodeCfg.Port)
	cli, err := http.NewRequest("POST", targetUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	cli.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(cli)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
