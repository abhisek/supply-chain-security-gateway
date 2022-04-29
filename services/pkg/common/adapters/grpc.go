package adapters

import (
	"log"
	"net"

	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
)

type GrpcAdapterConfigurer func(server *grpc.Server)

func GrpcStreamValidatorInterceptor() grpc.ServerOption {
	return grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_validator.StreamServerInterceptor(),
		),
	)
}

func GrpcUnaryValidatorInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_validator.UnaryServerInterceptor(),
		),
	)
}

func StartGrpcServer(name, host, port string, sopts []grpc.ServerOption, configure GrpcAdapterConfigurer) {
	addr := net.JoinHostPort(host, port)
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen on %s:%s - %s", host, port, err.Error())
	}

	server := grpc.NewServer(sopts...)
	configure(server)

	log.Printf("Starting %s gRPC server on %s:%s", name, host, port)
	err = server.Serve(listener)

	log.Fatalf("gRPC Server exit: %s", err.Error())
}
