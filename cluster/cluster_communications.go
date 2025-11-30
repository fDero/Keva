package cluster

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func internalHandlerFactory[T any](mutex *sync.Mutex, payload_callback func() any, handler func(T) int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request"})
			return
		}
		if mutex != nil {
			mutex.Lock()
			defer mutex.Unlock()
		}
		statusCode := handler(req)
		c.JSON(statusCode, payload_callback())
	}
}

func (rs *RaftSettings) StartClusterInternalServer(port string) error {
	api := gin.Default()
	api.POST("/v1/cluster/logsync", internalHandlerFactory(rs.mutex, rs.GetStateDescriptor, rs.handleLogSyncRequest))
	api.POST("/v1/cluster/votereq", internalHandlerFactory(rs.mutex, rs.GetStateDescriptor, rs.handleVoteRequest))
	api.POST("/v1/cluster/mkevent", internalHandlerFactory(rs.mutex, rs.GetStateDescriptor, rs.handleNewEventRequest))
	fmt.Println("Starting cluster internal server on port:", port)
	return api.Run(":" + port)
}

func (rs *RaftSettings) StartClusterEventLoop() {
	for {
		time.Sleep(rs.ping_frequency)
		rs.mutex.Lock()

		if rs.leader_identity == rs.self_identity {
			rs.DistributeLoggedEvents()
		} else if time.Since(rs.last_heartbeat) > rs.wait_time {
			rs.StartElection()
		}

		rs.mutex.Unlock()
	}
}
