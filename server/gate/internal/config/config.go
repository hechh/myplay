package config

import (
	"fmt"
	"myplay/common/pb"

	"github.com/hechh/framework"
	"github.com/hechh/framework/packet"
	"github.com/hechh/library/uerror"
	"github.com/hechh/library/yaml"
)

var (
	NodeCfg *yaml.NodeConfig
	GateCfg *GateConfig
)

type CommonConfig struct {
	yaml.CommonConfig
	TokenKey string `yaml:"token_key"`
	AesKey   string `yaml:"aes_key"`
}

type GateConfig struct {
	Mysql  map[int32]*yaml.DbConfig   `yaml:"mysql"`
	Redis  map[int32]*yaml.DbConfig   `yaml:"redis"`
	Etcd   *yaml.EtcdConfig           `yaml:"etcd"`
	Nats   *yaml.NatsConfig           `yaml:"nats"`
	Common *CommonConfig              `yaml:"common"`
	Server map[int32]*yaml.NodeConfig `yaml:"gate"`
}

func Load(cfg string, nodeId int32) error {
	item := &GateConfig{}
	if err := yaml.Load(cfg, item); err != nil {
		return err
	}
	nodeCfg, ok := item.Server[int32(nodeId)]
	if !ok {
		return uerror.New(-1, "节点配置(%d)不存在", nodeId)
	}
	GateCfg = item
	NodeCfg = nodeCfg

	// 初始化节点
	framework.Init(uint32(pb.NodeType_Gate), &packet.Node{
		Type: uint32(pb.NodeType_Gate),
		Id:   uint32(nodeId),
		Name: fmt.Sprintf("Gate%d", nodeId),
		Ip:   nodeCfg.Ip,
		Port: int32(nodeCfg.Port),
	})
	return nil
}
