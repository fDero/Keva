package misc

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"
)

type HostDesriptor struct {
	Address  string
	Port     string
	Identity string
}

func EncodeToBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func DecodeFromBase64(encoded string) string {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return ""
	}
	return string(decodedBytes)
}

func SplitN(s string, sep string, n int) []string {
	parts := strings.SplitN(s, sep, n)
	for len(parts) < n {
		parts = append(parts, "")
	}
	return parts[:n]
}

func AsJsonString(attribute string, value string) string {
	return "{ \"" + attribute + "\": \"" + value + "\" }"
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
