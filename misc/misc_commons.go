package misc

import (
	"encoding/base64"
	"net/http"
	"os"
)

type HostDesriptor struct {
	Address  string
	Port     string
	Identity string
}

func EncodeToBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func NormalizeEndpoint(endpoint string) string {
	if len(endpoint) > 0 && endpoint[0] == '/' {
		return endpoint[1:]
	}
	return endpoint
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CleanupConnection(resp *http.Response) {
	if resp != nil {
		resp.Body.Close()
	}
}
