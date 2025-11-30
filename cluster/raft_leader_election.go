package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fDero/keva/misc"
)

func (rs *RaftSettings) postElectionUpdate(votesReceived int) {
	votesNeeded := (rs.cluster_size / 2) + 1
	if votesReceived >= votesNeeded {
		rs.leader_identity = rs.self_identity
		rs.voted_for = ""
		fmt.Printf("Won election! I am now leader for epoch %d\n", rs.max_epoch)
	} else {
		fmt.Printf("Lost election. Got %d/%d votes\n", votesReceived, votesNeeded)
	}
}

func (rs *RaftSettings) createVoteRequest(peer ClusterNode) *http.Request {
	payload := map[string]any{
		"candidate_identity": rs.self_identity,
		"candidate_epoch":    rs.max_epoch,
		"candidate_last_log": rs.self_last_log,
	}
	json_data, _ := json.Marshal(payload)
	url := fmt.Sprintf("http://%s:%s/v1/cluster/votereq", peer.Address, peer.KevaPort)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func computeVotesIncrement(resp *http.Response, peer ClusterNode, err error) int {
	if err == nil && resp != nil && resp.StatusCode == RaftVoteGranted {
		fmt.Printf("Received vote from %s\n", peer.Identity)
		return 1
	}

	fmt.Printf("Failed to get vote from %s: %v\n", peer.Identity, err)
	return 0
}

func (rs *RaftSettings) StartElection() {
	rs.max_epoch += 1
	rs.voted_for = rs.self_identity
	votesReceived := 1

	fmt.Printf("Starting election for epoch %d\n", rs.max_epoch)

	for _, peer := range rs.other_nodes {
		req := rs.createVoteRequest(peer)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Do(req)
		votesReceived += computeVotesIncrement(resp, peer, err)
		defer misc.CleanupConnection(resp)
	}
	rs.postElectionUpdate(votesReceived)
}
