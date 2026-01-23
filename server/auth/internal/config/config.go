package config

import (
	"crypto/rand"
	"fmt"
	"myplay/common/pb"
	"myplay/common/table/global_config"

	"github.com/hechh/framework"
	"github.com/hechh/framework/packet"
	"github.com/hechh/library/myredis"
	"github.com/hechh/library/uerror"
	"github.com/hechh/library/yaml"
)

var (
	NodeCfg *yaml.NodeConfig
	AuthCfg *AuthConfig
)

type CommonConfig struct {
	Env         string `yaml:"env"`
	IsOpenPprof bool   `yaml:"is_open_pprof"`
	TokenKey    string `yaml:"token_key"`
	AesKey      string `yaml:"aes_key"`
}

type AuthConfig struct {
	Mysql  map[int32]*yaml.DbConfig   `yaml:"mysql"`
	Redis  map[int32]*yaml.DbConfig   `yaml:"redis"`
	Etcd   *yaml.EtcdConfig           `yaml:"etcd"`
	Nats   *yaml.NatsConfig           `yaml:"nats"`
	Common *CommonConfig              `yaml:"common"`
	Server map[int32]*yaml.NodeConfig `yaml:"auth"`
}

func Load(cfg string, nodeId int32) error {
	item := &AuthConfig{}
	if err := yaml.Load(cfg, item); err != nil {
		return err
	}
	nodeCfg, ok := item.Server[nodeId]
	if !ok {
		return uerror.New(-1, "节点配置(%d)不存在", nodeId)
	}
	AuthCfg = item
	NodeCfg = nodeCfg

	// 初始化节点
	framework.Init(uint32(pb.NodeType_Auth), &packet.Node{
		Type: uint32(pb.NodeType_Auth),
		Id:   uint32(nodeId),
		Name: fmt.Sprintf("Auth%d", nodeId),
		Ip:   nodeCfg.Ip,
		Port: int32(nodeCfg.Port),
	})
	return nil
}

func GetAesKey() string {
	return AuthCfg.Common.AesKey
}

func GetTokenKey() string {
	return AuthCfg.Common.TokenKey
}

func GenerateAes128Key() ([]byte, error) {
	key := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		return nil, uerror.Wrap(pb.ErrorCode_Unknown, err)
	}
	return key, nil
}

func GenerateUid() (uint64, error) {
	globalCfg := global_config.SGet(0)
	if globalCfg == nil {
		return 0, uerror.New(pb.ErrorCode_TableNotFound, "全局配置不存在")
	}
	client := myredis.Get("")
	val, err := client.Incr("user_register")
	if err != nil {
		return 0, uerror.Wrap(pb.ErrorCode_RedisFailed, err)
	}
	return globalCfg.UidBeginValue + uint64(val), nil
}
