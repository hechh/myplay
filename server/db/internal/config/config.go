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
	DbCfg   *DbConfig
)

type DbConfig struct {
	Mysql  map[int32]*yaml.DbConfig   `yaml:"mysql"`
	Redis  map[int32]*yaml.DbConfig   `yaml:"redis"`
	Etcd   *yaml.EtcdConfig           `yaml:"etcd"`
	Nats   *yaml.NatsConfig           `yaml:"nats"`
	Common *yaml.CommonConfig         `yaml:"common"`
	Server map[int32]*yaml.NodeConfig `yaml:"db"`
}

func Load(cfg string, nodeId int32) error {
	item := &DbConfig{}
	if err := yaml.Load(cfg, item); err != nil {
		return err
	}
	nodeCfg, ok := item.Server[int32(nodeId)]
	if !ok {
		return uerror.New(-1, "节点配置(%d)不存在", nodeId)
	}
	DbCfg = item
	NodeCfg = nodeCfg

	// 初始化节点
	framework.Init(uint32(pb.NodeType_Gate), &packet.Node{
		Type: uint32(pb.NodeType_Db),
		Id:   uint32(nodeId),
		Name: fmt.Sprintf("Db%d", nodeId),
		Ip:   nodeCfg.Ip,
		Port: int32(nodeCfg.Port),
	})

	return nil
}
