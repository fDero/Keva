package core

import (
	"net/http"

	"github.com/fDero/keva/misc"
	"github.com/gin-gonic/gin"
)

type Event struct {
	action    string
	key       string
	new_value string
}

func DecodeAndForward(processEvent func(Event) error) func(string) error {
	return func(encoded_event string) error {
		event := DecodeEvent(encoded_event)
		return processEvent(event)
	}
}

func DecodeEvent(encoded_event string) Event {
	parts := misc.SplitN(encoded_event, "|", 3)
	action := misc.DecodeFromBase64(parts[0])
	key := misc.DecodeFromBase64(parts[1])
	value := misc.DecodeFromBase64(parts[2])
	return Event{
		action:    action,
		key:       key,
		new_value: value,
	}
}

func RecieveRestEvent(action string, c *gin.Context) (Event, error) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read request body"})
		return Event{}, err
	}
	event := Event{
		action:    action,
		key:       c.Param("key"),
		new_value: string(body),
	}
	return event, nil
}

func (e Event) Encode() string {
	encoded_action := misc.EncodeToBase64(e.action)
	encoded_key := misc.EncodeToBase64(e.key)
	encoded_value := misc.EncodeToBase64(e.new_value)
	return encoded_action + "|" + encoded_key + "|" + encoded_value
}
