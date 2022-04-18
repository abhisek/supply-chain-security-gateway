package tap

import (
	"log"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	envoy_v3_ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type tapService struct {
	config *common_config.Config
}

func NewTapService(config *common_config.Config) (envoy_v3_ext_proc_pb.ExternalProcessorServer, error) {
	return &tapService{config: config}, nil
}

func (s *tapService) Process(srv envoy_v3_ext_proc_pb.ExternalProcessor_ProcessServer) error {
	log.Printf("Tap service: Handling stream")

	ctx := srv.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := srv.Recv()
		if err != nil {
			return status.Errorf(codes.Unknown, "Error receiving request: %v", err)
		}

		resp := &envoy_v3_ext_proc_pb.ProcessingResponse{}
		switch req.Request.(type) {
		case *envoy_v3_ext_proc_pb.ProcessingRequest_RequestHeaders:
			log.Printf("Handling request headers request")
			break
		case *envoy_v3_ext_proc_pb.ProcessingRequest_ResponseHeaders:
			log.Printf("Handling response headers request")
			break
		default:
			log.Printf("Unknown request type: %v", req.Request)
		}

		if err := srv.Send(resp); err != nil {
			log.Printf("Failed to send stream response: %v", err)
		}
	}
}
