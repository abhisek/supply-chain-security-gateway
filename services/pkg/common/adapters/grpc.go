package adapters

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type GrpcAdapterConfigurer func(server *grpc.Server)

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
