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
	GameCfg *GameConfig
)

type GameConfig struct {
	Mysql   map[int32]*yaml.DbConfig     `yaml:"mysql"`
	Redis   map[int32]*yaml.DbConfig     `yaml:"redis"`
	Nats    *yaml.NatsConfig             `yaml:"nats"`
	Router  *yaml.RouterConfig          `yaml:"router"`
	Cluster *yaml.ClusterConfig         `yaml:"cluster"`
	Common  *yaml.CommonConfig           `yaml:"common"`
	Server  map[int32]*yaml.NodeConfig   `yaml:"game"`
}

func Load(cfg string, nodeId int32) error {
	item := &GameConfig{}
	if err := yaml.Load(cfg, item); err != nil {
		return err
	}
	nodeCfg, ok := item.Server[int32(nodeId)]
	if !ok {
		return uerror.New(-1, "节点配置(%d)不存在", nodeId)
	}
	GameCfg = item
	NodeCfg = nodeCfg

	// 初始化节点
	framework.Init(uint32(pb.NodeType_Gate), &packet.Node{
		Type: uint32(pb.NodeType_Game),
		Id:   uint32(nodeId),
		Name: fmt.Sprintf("Game%d", nodeId),
		Ip:   nodeCfg.Ip,
		Port: int32(nodeCfg.Port),
	})
	return nil
}
