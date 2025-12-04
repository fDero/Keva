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

func InitializeSettings(
	other_nodes map[string]cluster.ClusterNode,
	self_config cluster.ClusterNode,
	history_file history.HistoryFile,
	storage *core.StorageSettings,
) (*cluster.RaftSettings, *core.ServerSettings) {
	global_mutex := &sync.Mutex{}
	rs := cluster.NewRaftSettings(
		other_nodes,
		self_config,
		global_mutex,
		history_file.GetLastEventID(),
		registerCallbackFactory(history_file, storage),
		history_file.GetEventByID,
	)
	ss := core.NewServerSettings(
		global_mutex,
		rs.GetLeaderDescriptor,
		storage.FetchRecord,
	)
	return &rs, &ss
}

func registerCallbackFactory(
	history_file history.HistoryFile,
	storage *core.StorageSettings,
) func(string) error {
	return func(event string) error {
		err := storage.ProcessEvent(event)
		if err != nil {
			return err
		}
		err = history_file.AppendEvent(event)
		fmt.Println("Appended event to history:", event)
		return err
	}
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

	storage := core.NewStorageSettings()
	rs, ss := InitializeSettings(other_nodes, self_config, history_file, storage)

	go rs.StartClusterEventLoop()
	go rs.StartClusterInternalServer(self_config.KevaPort)
	go ss.StartUserAPIServer(self_config.UserPort)

	select {}
}

func main() {
	err := App.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
