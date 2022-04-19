package tap

import (
	"context"
	"log"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	envoy_v3_ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type tapService struct {
	handlerChain TapHandlerChain
	config       *common_config.Config
}

func NewTapService(config *common_config.Config, registrations []TapHandlerRegistration) (envoy_v3_ext_proc_pb.ExternalProcessorServer, error) {
	return &tapService{config: config,
		handlerChain: TapHandlerChain{Handlers: registrations}}, nil
}

func (s *tapService) RegisterHandler(handler TapHandlerRegistration) {
	s.handlerChain.Handlers = append(s.handlerChain.Handlers, handler)
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
			err = s.handleRequestHeaders(ctx,
				req.Request.(*envoy_v3_ext_proc_pb.ProcessingRequest_RequestHeaders),
				resp)
			break
		case *envoy_v3_ext_proc_pb.ProcessingRequest_ResponseHeaders:
			err = s.handleResponseHeaders(ctx,
				req.Request.(*envoy_v3_ext_proc_pb.ProcessingRequest_ResponseHeaders),
				resp)
			break
		default:
			log.Printf("Unknown request type: %v", req.Request)
		}

		// TODO: How should we handle this behavior?
		if err != nil {
			log.Printf("Error in handling processing req: %v", err)
		}

		if err := srv.Send(resp); err != nil {
			log.Printf("Failed to send stream response: %v", err)
		}
	}
}

func (s *tapService) handleRequestHeaders(ctx context.Context,
	req *envoy_v3_ext_proc_pb.ProcessingRequest_RequestHeaders,
	resp *envoy_v3_ext_proc_pb.ProcessingResponse) error {
	for _, registration := range s.handlerChain.Handlers {
		err := registration.Handler.HandleRequestHeaders(ctx, req, resp)
		if !registration.ContinueOnError && err != nil {
			log.Printf("Unable to continue on tap handler error: %v", err)
			return err
		}
	}

	return nil
}

func (s *tapService) handleResponseHeaders(ctx context.Context,
	req *envoy_v3_ext_proc_pb.ProcessingRequest_ResponseHeaders,
	resp *envoy_v3_ext_proc_pb.ProcessingResponse) error {
	for _, registration := range s.handlerChain.Handlers {
		err := registration.Handler.HandleResponseHeaders(ctx, req, resp)
		if !registration.ContinueOnError && err != nil {
			log.Printf("Unable to continue on tap handler error: %v", err)
			return err
		}
	}

	return nil
}
