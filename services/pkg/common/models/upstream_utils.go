package models

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

var (
	errIncorrectPrefix           = errors.New("incorrect path prefix")
	errIncorrectMaven2Path       = errors.New("incorrect maven2 path")
	errUnimplementedUpstreamType = errors.New("path resolver for upstream type is not implemented")
)

func GetArtefactByHostAndPath(upstreams []ArtefactUpStream, host, path string) (Artefact, error) {
	for _, upstream := range upstreams {
		if upstream.MatchHost(host) && upstream.MatchPath(path) {
			return upstream.Path2Artefact(path)
		}
	}

	return Artefact{}, fmt.Errorf("no artefact resolved using %s/%s", host, path)
}

func (s ArtefactUpStream) MatchHost(host string) bool {
	return (s.RoutingRule.Host == host)
}

func (s ArtefactUpStream) MatchPath(path string) bool {
	path = filepath.Clean(path)
	return strings.HasPrefix(path, s.RoutingRule.Prefix)
}

// Resolve an HTTP request path for this artefact into an Arteface model
func (s ArtefactUpStream) Path2Artefact(path string) (Artefact, error) {
	path = filepath.Clean(path)
	if !strings.HasPrefix(path, s.RoutingRule.Prefix) {
		return Artefact{}, errIncorrectPrefix
	}

	path = strings.TrimPrefix(path, s.RoutingRule.Prefix)
	if path != "" && path[0] == '/' {
		path = path[1:]
	}

	parts := strings.Split(path, "/")
	switch s.Type {
	case ArtefactSourceTypeMaven2:
		return artefactForMaven2(parts)
	default:
		return Artefact{}, errUnimplementedUpstreamType
	}
}

func artefactForMaven2(parts []string) (Artefact, error) {
	if len(parts) < 4 {
		return Artefact{}, errIncorrectMaven2Path
	}

	// Ignore the filename
	_ = parts[:len(parts)-1]

	version := parts[len(parts)-2]
	name := parts[len(parts)-3]

	parts = parts[:len(parts)-3]
	group := strings.Join(parts, ".")

	return NewArtefact(ArtefactSource{Type: ArtefactSourceTypeMaven2},
		name, group, version), nil
}
