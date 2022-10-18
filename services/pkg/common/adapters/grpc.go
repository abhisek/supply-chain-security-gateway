package adapters

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	grpcotel "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

type GrpcAdapterConfigurer func(server *grpc.Server)
type GrpcClientConfigurer func(conn *grpc.ClientConn)

var (
	NoGrpcDialOptions = []grpc.DialOption{}
	NoGrpcConfigurer  = func(conn *grpc.ClientConn) {}
)

func StartGrpcMtlsServer(name, serverName, host, port string, sopts []grpc.ServerOption, configure GrpcAdapterConfigurer) {
	tc, err := utils.TlsConfigFromEnvironment(serverName)
	if err != nil {
		log.Fatalf("Failed to setup TLS from environment: %v", err)
	}

	creds := credentials.NewTLS(&tc)
	sopts = append(sopts, grpc.Creds(creds))

	StartGrpcServer(name, host, port, sopts, configure)
}

func StartGrpcServer(name, host, port string, sopts []grpc.ServerOption, configure GrpcAdapterConfigurer) {
	addr := net.JoinHostPort(host, port)
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen on %s:%s - %s", host, port, err.Error())
	}

	sopts = append(sopts, grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpcotel.UnaryServerInterceptor(),
			grpc_validator.UnaryServerInterceptor(),
		),
	))

	sopts = append(sopts, grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpcotel.StreamServerInterceptor(),
			grpc_validator.StreamServerInterceptor(),
		),
	))

	server := grpc.NewServer(sopts...)
	configure(server)

	log.Printf("Starting %s gRPC server on %s:%s", name, host, port)
	err = server.Serve(listener)

	log.Fatalf("gRPC Server exit: %s", err.Error())
}

func GrpcMtlsClient(name, serverName, host, port string, dopts []grpc.DialOption, configurer GrpcClientConfigurer) (*grpc.ClientConn, error) {
	tc, err := grpcTransportCredentials(serverName)
	if err != nil {
		return nil, fmt.Errorf("failed to setup client transport credentials: %w", err)
	}

	dopts = append(dopts, tc)
	return grpcClient(name, host, port, dopts, configurer)
}

func GrpcInsecureClient(name, host, port string, dopts []grpc.DialOption, configurer GrpcClientConfigurer) (*grpc.ClientConn, error) {
	tc := grpc.WithTransportCredentials(insecure.NewCredentials())
	dopts = append(dopts, tc)
	return grpcClient(name, host, port, dopts, configurer)
}

func grpcClient(name, host, port string, dopts []grpc.DialOption, configurer GrpcClientConfigurer) (*grpc.ClientConn, error) {
	log.Printf("[%s] Connecting to gRPC server %s:%s", name, host, port)

	dopts = append(dopts, grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	dopts = append(dopts, grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))

	retry := 5
	t := 1
	conn, err := grpc.Dial(net.JoinHostPort(host, port), dopts...)
	for err != nil && t < retry {
		log.Printf("[%d/%d] Retrying due to failure: %v", t, retry, err)
		conn, err = grpc.Dial(net.JoinHostPort(host, port), dopts...)

		time.Sleep(1 * time.Second)
		t += 1
	}

	if err != nil {
		return nil, err
	}

	configurer(conn)
	return conn, nil
}

func grpcTransportCredentials(serverName string) (grpc.DialOption, error) {
	tlsConfig, err := utils.TlsConfigFromEnvironment(serverName)
	if err != nil {
		return nil, err
	}

	creds := credentials.NewTLS(&tlsConfig)
	return grpc.WithTransportCredentials(creds), nil
}
