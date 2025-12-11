package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fDero/keva/misc"
	"github.com/gin-gonic/gin"
)

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

func (ss *ServerSettings) discoverLeader() misc.HostDesriptor {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	return ss.getLeaderInfoCallback()
}

func (ss *ServerSettings) fetchRecord(ev Event) (int, string) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	value, present := ss.fetchElementCallback(ev.key)
	if !present {
		fmt.Printf("Cannot find record: key=%s\n", ev.key)
		return http.StatusNoContent, `{"value": ""}`
	}
	fmt.Printf("Fetching record: key=%s, value=%s\n", ev.key, value)
	return http.StatusOK, value
}

func (ss *ServerSettings) forwardEvent(ev Event) (int, string) {
	for {
		status := ss.forwardToLeader(
			"/v1/cluster/mkevent/",
			gin.H{"new_event": ev.Encode()},
		)
		if status >= 200 && status < 300 {
			return http.StatusOK, misc.AsJsonString("done", "true")
		}
	}
}
