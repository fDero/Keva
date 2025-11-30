package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fDero/keva/misc"
)

func (rs *RaftSettings) makeSyncBasicRequest(log int64, peer ClusterNode) *http.Request {
	payload := map[string]any{
		"leader_identity":   rs.self_identity,
		"leader_epoch":      rs.max_epoch,
		"leader_last_event": rs.restore_event_callback(log),
		"leader_last_log":   log,
	}
	json_data, _ := json.Marshal(payload)
	url := fmt.Sprintf("http://%s:%s/v1/cluster/logsync", peer.Address, peer.KevaPort)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func (rs *RaftSettings) makeSyncRecoveryRequest(resp *http.Response, peer ClusterNode) *http.Request {
	var json_resp_structure struct {
		FollowerIdentity string `json:"identity"`
		FollowerLastLog  int64  `json:"last_log"`
		FollowerEpoch    int64  `json:"epoch"`
		FollowerLeader   string `json:"leader_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&json_resp_structure); err != nil {
		log.Printf("Failed to decode JSON: %v", err)
	}
	var log_id int64 = json_resp_structure.FollowerLastLog + 1
	return rs.makeSyncBasicRequest(log_id, peer)
}

func (rs *RaftSettings) DistributeLoggedEvents() {
	fmt.Printf("=== LEADER distributing logs. My identity: '%s' ===\n", rs.self_identity)
	for _, peer := range rs.other_nodes {
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Do(rs.makeSyncBasicRequest(rs.self_last_log, peer))
		defer misc.CleanupConnection(resp)
		if err != nil {
			fmt.Printf("Failed to send heartbeat to %s: %v\n", peer.Identity, err)
		}
		if resp != nil && resp.StatusCode == RaftFollowerLaggingBehind {
			resp2, _ := client.Do(rs.makeSyncRecoveryRequest(resp, peer))
			defer misc.CleanupConnection(resp2)
		}
	}
}
