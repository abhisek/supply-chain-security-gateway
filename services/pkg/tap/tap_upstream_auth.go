package tap

import (
	"log"

	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
	envoy_v3_ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

func (s *tapService) applyUpstreamAuth(req *envoy_v3_ext_proc_pb.ProcessingRequest_RequestHeaders,
	resp *envoy_v3_ext_proc_pb.ProcessingResponse_RequestHeaders) error {

	path, err := findHeaderValue(req, "path")
	if err != nil {
		return err
	}

	upstream, err := common_models.GetUpstreamByHostAndPath(s.config.Global.Upstreams,
		"", path)
	if err != nil {
		return err
	}

	if !upstream.NeedUpstreamAuthentication() {
		log.Printf("Upstream %s do not need authentication", upstream.Name)
		return nil
	}

	return nil
}
