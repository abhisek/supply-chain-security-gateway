package tap

import (
	"fmt"
	"strings"

	envoy_v3_ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

// Header keys are stored as ":key" by envoy
func findHeaderValue(req *envoy_v3_ext_proc_pb.ProcessingRequest_RequestHeaders,
	key string) (string, error) {
	for _, h := range req.RequestHeaders.Headers.Headers {
		if strings.EqualFold(":"+key, h.Key) {
			return h.Value, nil
		}
	}

	return "", fmt.Errorf("header with key: %s not found", key)
}

func findHostAndPath(req *envoy_v3_ext_proc_pb.ProcessingRequest_RequestHeaders) (string, string, error) {
	path, err := findHeaderValue(req, "path")
	if err != nil {
		return "", "", fmt.Errorf("failed to find path in req: %w", err)
	}

	host, err := findHeaderValue(req, "host")
	if err != nil {
		return "", "", fmt.Errorf("failed to find host in req: %w", err)
	}

	return host, path, nil
}
