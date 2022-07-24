package route

import "fmt"

type routeHandler struct {
	pathPattern string
}

type routeMatch struct {
	is_match bool
	labels   map[string]string
}

func NewRouteHandler(pathP string) (*routeHandler, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (h *routeHandler) Match(path string) routeMatch {
	return routeMatch{}
}

func (r *routeMatch) IsMatch() bool {
	return r.is_match
}

func (r *routeMatch) Labels() map[string]string {
	return r.labels
}
