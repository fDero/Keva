package core

import (
	"fmt"
	"sync"

	"github.com/fDero/keva/cluster"
	"github.com/fDero/keva/history"
	"github.com/fDero/keva/server"
	"github.com/fDero/keva/storage"
)

func InitializeSettings(
	other_nodes map[string]cluster.ClusterNode,
	self_config cluster.ClusterNode,
	history_file history.HistoryFile,
	storage *storage.StorageSettings,
) (*cluster.RaftSettings, *server.ServerSettings) {
	global_mutex := &sync.Mutex{}
	rs := cluster.NewRaftSettings(
		other_nodes,
		self_config,
		global_mutex,
		history_file.GetLastEventID(),
		registerCallbackFactory(history_file, storage),
		history_file.GetEventByID,
	)
	ss := server.NewServerSettings(
		global_mutex,
		rs.GetLeaderDescriptor,
		storage.FetchRecord,
	)
	return &rs, &ss
}

func registerCallbackFactory(
	history_file history.HistoryFile,
	storage *storage.StorageSettings,
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
