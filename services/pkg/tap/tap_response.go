package tap

import envoy_v3_ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"

type tapResponse struct {
	pr *envoy_v3_ext_proc_pb.ProcessingResponse
}

func buildTapResponse() *tapResponse {
	return &tapResponse{
		pr: &envoy_v3_ext_proc_pb.ProcessingResponse{},
	}
}

func (r *tapResponse) Response() *envoy_v3_ext_proc_pb.ProcessingResponse {
	return r.pr
}
