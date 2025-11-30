package cluster

import (
	"sync"
	"time"

	"github.com/fDero/keva/misc"
)

const (
	RaftVoteGranted           = 240
	RaftSyncAccepted          = 250
	RaftEventWritten          = 260
	RaftVoteDenied            = 441
	RaftFollowerLaggingBehind = 451
	RaftLeaderLaggingBehind   = 452
	RaftEventIgnored          = 561
)

type ClusterNode struct {
	Identity string `toml:"identity"`
	Address  string `toml:"address"`
	KevaPort string `toml:"keva_port"`
	UserPort string `toml:"user_port"`
	WaitTime string `toml:"wait_time"`
}

type RaftSettings struct {
	self_host              misc.HostDesriptor
	self_identity          string
	max_epoch              int64
	self_epoch             int64
	self_last_log          int64
	leader_identity        string
	voted_for              string
	cluster_size           int
	last_heartbeat         time.Time
	wait_time              time.Duration
	ping_frequency         time.Duration
	mutex                  *sync.Mutex
	save_event_callback    func(event string) error
	restore_event_callback func(id int64) string
	other_nodes            map[string]ClusterNode
}
