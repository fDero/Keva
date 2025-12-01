package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/fDero/keva/misc"
	"github.com/gin-gonic/gin"
)

type ServerSettings struct {
	getLeaderInfoCallback func() misc.HostDesriptor
	fetchElementCallback  func(key string) (string, bool)
	mutex                 *sync.Mutex
}

func NewServerSettings(
	mutex *sync.Mutex,
	getLeaderInfoCallback func() misc.HostDesriptor,
	fetchElementCallback func(key string) (string, bool),
) ServerSettings {
	return ServerSettings{
		getLeaderInfoCallback: getLeaderInfoCallback,
		fetchElementCallback:  fetchElementCallback,
		mutex:                 mutex,
	}
}

func forwarderFactory(handler func(string, string) int) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		var jsonBody struct {
			Value string `json:"value"`
		}
		if err := c.ShouldBindJSON(&jsonBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON or missing value field"})
			return
		}
		statusCode := handler(key, jsonBody.Value)
		c.JSON(statusCode, gin.H{})
	}
}

func handlerFactory(handler func(string) (int, string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		status_code, json_string := handler(key)
		json_obj := make(map[string]any)
		err := json.Unmarshal([]byte(json_string), &json_obj)
		if err != nil {
			c.JSON(http.StatusInternalServerError, json_string)
			return
		}
		c.JSON(status_code, json_obj)
	}
}

func (ss *ServerSettings) StartUserAPIServer(port string) error {
	api := gin.Default()
	api.POST("/v1/storage/key/:key", forwarderFactory(ss.upsertRecord))
	api.DELETE("/v1/storage/key/:key", forwarderFactory(ss.deleteRecord))
	api.GET("/v1/storage/key/:key", handlerFactory(ss.fetchRecord))
	fmt.Println("Starting storage user API server on port:", port)
	return api.Run(":" + port)
}
