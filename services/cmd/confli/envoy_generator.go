package main

import (
	"fmt"
	"os"

	"github.com/abhisek/supply-chain-gateway/services/gen"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
	"google.golang.org/protobuf/types/known/structpb"

	envoy_bootstrap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	envoy_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/golang/protobuf/jsonpb"
)

/**
Envoy config references:

https://www.envoyproxy.io/docs/envoy/latest/api-v3/api
https://www.envoyproxy.io/docs/envoy/latest/configuration/overview/examples
*/

type envoyConfigGenerator struct {
	configFile string
}

func newEnvoyConfigGenerator(path string) *envoyConfigGenerator {
	return &envoyConfigGenerator{
		configFile: path,
	}
}

func (g *envoyConfigGenerator) generate() error {
	repo, err := config.NewConfigFileRepository(g.configFile, false, false)
	if err != nil {
		return fmt.Errorf("failed to create config repo: %w", err)
	}

	cfg, err := repo.LoadGatewayConfiguration()
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	err = cfg.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	bootstrapConfig, err := apiGenerateEnvoyConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to generate envoy config: %w", err)
	}

	return printEnvoyBootstrapConfig(bootstrapConfig)
}

func printEnvoyBootstrapConfig(cfg *envoy_bootstrap_v3.Bootstrap) error {
	m := jsonpb.Marshaler{Indent: "  "}
	data, err := m.MarshalToString(cfg)

	if err != nil {
		return err
	}

	_, err = os.Stdout.Write([]byte(data))
	return err
}

func apiGenerateEnvoyConfig(gateway *gen.GatewayConfiguration) (*envoy_bootstrap_v3.Bootstrap, error) {
	bootstrap := envoy_bootstrap_v3.Bootstrap{
		Node: &envoy_core_v3.Node{
			Id:      envoyNodeId(gateway.Info.Id),
			Cluster: gateway.Info.Name,
			Metadata: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"domain": {
						Kind: &structpb.Value_StringValue{StringValue: gateway.Info.Domain},
					},
				},
			},
		},
		StaticResources: &envoy_bootstrap_v3.Bootstrap_StaticResources{
			Listeners: make([]*envoy_listener_v3.Listener, 0),
			Clusters:  make([]*envoy_cluster_v3.Cluster, 0),
		},
	}

	return &bootstrap, nil
}

func envoyNodeId(gid string) string {
	return fmt.Sprintf("%s--%s", gid, utils.NewUniqueId())
}
