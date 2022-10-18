package tap

import (
	"reflect"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/logger"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_v3_ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

const (
	tapSignatureHeaderName  = "x-gateway-tap"
	tapSignatureHeaderValue = "true"
)

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

func (r *tapResponse) WithProcessingResponseRequestHeaders() *tapResponse {
	if r.pr.Response != nil {
		logger.Errorf("this tap response has a envoy processing response already set: %v",
			reflect.TypeOf(r.pr.Response))
		return r
	}

	r.pr.Response = &envoy_v3_ext_proc_pb.ProcessingResponse_RequestHeaders{
		RequestHeaders: &envoy_v3_ext_proc_pb.HeadersResponse{
			Response: &envoy_v3_ext_proc_pb.CommonResponse{
				Status: envoy_v3_ext_proc_pb.CommonResponse_CONTINUE,
			},
		},
	}

	return r
}

func (r *tapResponse) WithTapSignature() *tapResponse {
	if r.pr.Response == nil {
		logger.Errorf("this tap response is not initialized")
		return r
	}

	return r
}

func (r *tapResponse) SetResponseHeader(key, value string) *tapResponse {
	if res, ok := r.AsProcessingResponseResponseHeaders(); ok {
		res.ResponseHeaders.Response.HeaderMutation.SetHeaders = append(res.ResponseHeaders.Response.HeaderMutation.SetHeaders,
			&envoy_config_core_v3.HeaderValueOption{
				Header: &envoy_config_core_v3.HeaderValue{
					Key:   tapSignatureHeaderName,
					Value: tapSignatureHeaderValue,
				},
			})
	}

	return r
}

func (r *tapResponse) AsProcessingResponseResponseHeaders() (*envoy_v3_ext_proc_pb.ProcessingResponse_ResponseHeaders, bool) {
	cr, ok := r.pr.Response.(*envoy_v3_ext_proc_pb.ProcessingResponse_ResponseHeaders)
	return cr, ok
}

func (r *tapResponse) AsProcessingResponseRequestHeaders() (*envoy_v3_ext_proc_pb.ProcessingResponse_RequestHeaders, bool) {
	cr, ok := r.pr.Response.(*envoy_v3_ext_proc_pb.ProcessingResponse_RequestHeaders)
	return cr, ok
}

func (r *tapResponse) WithProcessingResponseResponseHeaders() *tapResponse {
	if r.pr.Response != nil {
		logger.Errorf("this tap response has a envoy processing response already set: %v",
			reflect.TypeOf(r.pr.Response))
		return r
	}

	r.pr.Response = &envoy_v3_ext_proc_pb.ProcessingResponse_ResponseHeaders{
		ResponseHeaders: &envoy_v3_ext_proc_pb.HeadersResponse{
			Response: &envoy_v3_ext_proc_pb.CommonResponse{
				HeaderMutation: &envoy_v3_ext_proc_pb.HeaderMutation{
					SetHeaders: []*envoy_config_core_v3.HeaderValueOption{},
				},
			},
		},
	}

	return r
}
