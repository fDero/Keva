package cluster

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fDero/keva/misc"
)

func LoadClusterConfig(config_filepath string) ([]ClusterNode, error) {
	type TomlNodes struct {
		Nodes []ClusterNode `toml:"node"`
	}
	var config TomlNodes
	if _, err := toml.DecodeFile(config_filepath, &config); err != nil {
		return nil, fmt.Errorf("failed to parse TOML cluster configuration file: %w", err)
	}
	return config.Nodes, nil
}

func SplitClusterNodes(self_identity string, all_nodes []ClusterNode) (ClusterNode, map[string]ClusterNode) {
	other_nodes := make(map[string]ClusterNode)
	var self_config *ClusterNode = nil

	for _, node := range all_nodes {
		if node.Identity == self_identity {
			self_config = &node
		} else {
			other_nodes[node.Identity] = node
		}
	}
	return *self_config, other_nodes
}

func NewRaftSettings(
	other_nodes map[string]ClusterNode,
	self_config ClusterNode,
	mutex *sync.Mutex,
	last_log int64,
	inputsave_event_callback func(string) error,
	inputrestore_event_callback func(int64) string,

) RaftSettings {
	wait_time, _ := strconv.Atoi(self_config.WaitTime)
	return RaftSettings{
		self_host:              self_config.AsHostDescriptor(),
		self_identity:          self_config.Identity,
		max_epoch:              0,
		self_epoch:             0,
		self_last_log:          last_log,
		leader_identity:        "",
		voted_for:              "",
		cluster_size:           len(other_nodes) + 1,
		last_heartbeat:         time.Now(),
		wait_time:              time.Second * time.Duration(wait_time),
		ping_frequency:         time.Second * 1,
		mutex:                  mutex,
		save_event_callback:    inputsave_event_callback,
		restore_event_callback: inputrestore_event_callback,
		other_nodes:            other_nodes,
	}
}

func (rs *RaftSettings) GetStateDescriptor() any {
	return map[string]any{
		"identity":  rs.self_identity,
		"last_log":  rs.self_last_log,
		"epoch":     rs.max_epoch,
		"leader_id": rs.leader_identity,
	}
}

func (rs *RaftSettings) GetLeaderDescriptor() misc.HostDesriptor {
	if leader, exists := rs.other_nodes[rs.leader_identity]; exists {
		return misc.HostDesriptor{
			Identity: leader.Identity,
			Address:  leader.Address,
			Port:     leader.KevaPort,
		}
	}
	fmt.Printf("Leader %s not found among other nodes, assuming self is leader\n", rs.leader_identity)
	fmt.Printf("Self identity: %s\n <--\n", rs.self_identity)
	return rs.self_host
}

func (cn *ClusterNode) AsHostDescriptor() misc.HostDesriptor {
	return misc.HostDesriptor{
		Identity: cn.Identity,
		Address:  cn.Address,
		Port:     cn.KevaPort,
	}
}
