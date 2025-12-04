package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fDero/keva/misc"
	"github.com/gin-gonic/gin"
)

func (ss *ServerSettings) discoverLeader() misc.HostDesriptor {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	return ss.getLeaderInfoCallback()
}

func (ss *ServerSettings) forwardToLeader(endpoint string, body map[string]any) int {
	leader := ss.discoverLeader()
	normalized_endpoint := misc.NormalizeEndpoint(endpoint)
	url := fmt.Sprintf("http://%s:%s/%s", leader.Address, leader.Port, normalized_endpoint)
	json_body, _ := json.Marshal(body)
	fmt.Printf("Forwarding to leader: %s at %s:%s\n", leader.Identity, leader.Address, leader.Port)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(json_body))
	if err != nil {
		fmt.Printf("Error forwarding to leader %s: %v\n", leader.Identity, err)
		return http.StatusServiceUnavailable
	}
	return resp.StatusCode
}

func (ss *ServerSettings) fetchRecord(key string) (int, string) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	value, present := ss.fetchElementCallback(key)
	if !present {
		fmt.Printf("Cannot find record: key=%s\n", key)
		return http.StatusNoContent, value
	}
	fmt.Printf("Fetching record: key=%s, value=%s\n", key, value)
	return http.StatusOK, value
}

func (ss *ServerSettings) upsertRecord(key string, value string) int {
	for {
		key_b64 := misc.EncodeToBase64(key)
		value_b64 := misc.EncodeToBase64(value)
		event := fmt.Sprintf("UPSERT|%s|%s", key_b64, value_b64)
		status := ss.forwardToLeader(
			"/v1/cluster/mkevent/",
			gin.H{"new_event": event},
		)
		if status >= 200 && status < 300 {
			return http.StatusOK
		}
	}
}

func (ss *ServerSettings) deleteRecord(key string, value string) int {
	for {
		key_b64 := misc.EncodeToBase64(key)
		event := fmt.Sprintf("DELETE|%s", key_b64)
		status := ss.forwardToLeader(
			"/v1/cluster/mkevent/",
			gin.H{"new_event": event},
		)
		if status >= 200 && status < 300 {
			return http.StatusOK
		}
	}
}
