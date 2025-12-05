package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fDero/keva/cluster"
	"github.com/fDero/keva/core"
	"github.com/fDero/keva/history"
	"github.com/fDero/keva/misc"
	"github.com/urfave/cli/v2"
)

var App = &cli.App{
	Name:  "keva",
	Usage: "A distributed fault-tolerant key-value storage system",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "config-file",
			Usage:    "TOML encoded file with cluster configuration",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "working-directory",
			Usage:    "Directory where internal system state can be stored for persistence",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "node-identity",
			Usage:    "Name that uniquely identifies this node in the cluster",
			Value:    misc.FailOnError(os.Hostname()),
			Required: false,
		},
	},
	Action: executeCommand,
}

func executeCommand(c *cli.Context) error {
	var cluster_nodes []cluster.ClusterNode
	var errors [5]error

	config_file := c.String("config-file")
	workidr := c.String("working-directory")
	self_identity := c.String("node-identity")

	lockfile := filepath.Join(workidr, "lock.json")
	errors[0] = misc.PurgeUnusedLockFiles(lockfile)
	errors[1] = misc.CreateLockFile(lockfile)

	cluster_nodes, errors[2] = cluster.LoadClusterConfig(config_file)
	self_config, other_nodes := cluster.SplitClusterNodes(self_identity, cluster_nodes)

	if err := misc.FirstOfManyErrorsOrNone(errors[:3]); err != nil {
		return fmt.Errorf("failed to initialize cluster node: %w", err)
	}

	var history_file history.HistoryFile
	pm := misc.NewDiskPersistenceHandler(workidr, "history.dat")

	default_fallback_header := history.GetDefaultHistoryFileHeader(0, 0)
	history_file, errors[4] = history.NewHistoryFile(pm, default_fallback_header)

	if err := misc.FirstOfManyErrorsOrNone(errors[:5]); err != nil {
		return fmt.Errorf("failed to initialize cluster node: %w", err)
	}

	global_mutex := &sync.Mutex{}
	storage := core.NewStorageSettings()

	raft_settings := cluster.NewRaftSettings(
		other_nodes,
		self_config,
		global_mutex,
		history_file.GetLastEventID(),
		misc.ProcessingPipeline(history_file.AppendEvent, storage.ProcessEvent),
		history_file.GetEventByID,
	)

	server_settings := core.NewServerSettings(
		global_mutex,
		raft_settings.GetLeaderDescriptor,
		storage.FetchRecord,
	)

	go raft_settings.StartClusterEventLoop()
	go raft_settings.StartClusterInternalServer(self_config.KevaPort)
	go server_settings.StartUserAPIServer(self_config.UserPort)

	select {}
}

func main() {
	err := App.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
