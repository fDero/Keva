package core

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

func handlerFactory(action string, handler func(Event) (int, string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		event, err := RecieveRestEvent(action, c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read request body"})
			return
		}
		status_code, json_string := handler(event)
		json_obj := make(map[string]any)
		err = json.Unmarshal([]byte(json_string), &json_obj)
		if err != nil {
			c.JSON(http.StatusInternalServerError, json_string)
			return
		}
		c.JSON(status_code, json_obj)
	}
}

func (ss *ServerSettings) StartUserAPIServer(port string) error {
	api := gin.Default()
	api.POST("/v1/storage/key/:key", handlerFactory("UPSERT/VALUE", ss.forwardEvent))
	api.DELETE("/v1/storage/key/:key", handlerFactory("DELETE/VALUE", ss.forwardEvent))
	api.GET("/v1/storage/key/:key", handlerFactory("READ/VALUE", ss.fetchRecord))
	fmt.Println("Starting storage user API server on port:", port)
	return api.Run(":" + port)
}
