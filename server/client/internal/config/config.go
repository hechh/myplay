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
	NodeCfg   *yaml.NodeConfig
	ClientCfg *ClientConfig
)

type ClientConfig struct {
	Mysql   map[int32]*yaml.DbConfig     `yaml:"mysql"`
	Redis   map[int32]*yaml.DbConfig     `yaml:"redis"`
	Nats    *yaml.NatsConfig             `yaml:"nats"`
	Router  *yaml.RouterConfig          `yaml:"router"`
	Cluster *yaml.ClusterConfig         `yaml:"cluster"`
	Common  *yaml.CommonConfig           `yaml:"common"`
	Server  map[int32]*yaml.NodeConfig   `yaml:"client"`
	Gates   map[int32]*yaml.NodeConfig   `yaml:"db"`
}

func Load(cfg string, nodeId int32) error {
	item := &ClientConfig{}
	if err := yaml.Load(cfg, item); err != nil {
		return err
	}
	nodeCfg, ok := item.Server[int32(nodeId)]
	if !ok {
		return uerror.New(-1, "节点配置(%d)不存在", nodeId)
	}
	ClientCfg = item
	NodeCfg = nodeCfg

	// 初始化节点
	framework.Init(uint32(pb.NodeType_Gate), &packet.Node{
		Type: uint32(pb.NodeType_Client),
		Id:   uint32(nodeId),
		Name: fmt.Sprintf("Client%d", nodeId),
		Ip:   nodeCfg.Ip,
		Port: int32(nodeCfg.Port),
	})

	return nil
}

func GetWsUrl(nodeId uint32) (string, error) {
	cfg, ok := ClientCfg.Gates[int32(nodeId)]
	if !ok {
		return "", uerror.Err(pb.ErrorCode_Unknown, "gate节点(%d)配置不存在", nodeId)
	}
	return fmt.Sprintf("ws://%s:%d/ws", cfg.Ip, cfg.Port), nil
}
