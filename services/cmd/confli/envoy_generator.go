package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/abhisek/supply-chain-gateway/services/gen"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"

	envoy_bootstrap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	envoy_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoy_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_extension_http_connection_manager_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	envoy_extension_tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	envoy_extension_http_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

const (
	pdpSvcName  = "ext-authz-pdp"
	pdpHostName = "pdp"
	PdpPort     = "9001"

	tapSvcName  = "ext-proc-tap"
	tapHostName = "tap"
	tapPort     = "9001"
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

	listener, err := envoyGenerateStaticListener(gateway)
	if err != nil {
		return nil, fmt.Errorf("failed to generated static listener: %w", err)
	}

	clusters, err := envoyGenerateStaticClusters(gateway)
	if err != nil {
		return nil, fmt.Errorf("failed to generate static clusters: %w", err)
	}

	bootstrap.StaticResources.Listeners = append(bootstrap.StaticResources.Listeners, listener)
	bootstrap.StaticResources.Clusters = append(bootstrap.StaticResources.Clusters, clusters...)

	return &bootstrap, nil
}

func envoyGenerateStaticListener(gateway *gen.GatewayConfiguration) (*envoy_listener_v3.Listener, error) {
	listener := &envoy_listener_v3.Listener{
		Address: &envoy_core_v3.Address{
			Address: &envoy_core_v3.Address_SocketAddress{
				SocketAddress: &envoy_core_v3.SocketAddress{
					Address:       gateway.Listener.Host,
					PortSpecifier: &envoy_core_v3.SocketAddress_PortValue{PortValue: gateway.Listener.Port},
				},
			},
		},
		FilterChains: make([]*envoy_listener_v3.FilterChain, 0),
	}

	routeConfig := &envoy_route_v3.RouteConfiguration{
		Name:                 "local_route",
		VirtualHosts:         make([]*envoy_route_v3.VirtualHost, 0),
		ResponseHeadersToAdd: make([]*envoy_core_v3.HeaderValueOption, 0),
		RequestHeadersToAdd:  make([]*envoy_core_v3.HeaderValueOption, 0),
	}

	http_connection_manager := &envoy_extension_http_connection_manager_v3.HttpConnectionManager{
		StatPrefix: "ingress_http",
		RouteSpecifier: &envoy_extension_http_connection_manager_v3.HttpConnectionManager_RouteConfig{
			RouteConfig: routeConfig,
		},
	}

	vhosts := &envoy_route_v3.VirtualHost{
		Name:    "catch_all_vhost",
		Domains: []string{"*"},
		Routes:  make([]*envoy_route_v3.Route, 0),
	}

	for _, upstream := range gateway.Upstreams {
		route := &envoy_route_v3.Route{
			Match: &envoy_route_v3.RouteMatch{
				PathSpecifier: &envoy_route_v3.RouteMatch_Prefix{
					Prefix: upstream.Route.PathPrefix,
				},
			},
			Action: &envoy_route_v3.Route_Route{
				Route: &envoy_route_v3.RouteAction{
					Timeout: &durationpb.Duration{
						Seconds: 0,
						Nanos:   0,
					},
					ClusterSpecifier: &envoy_route_v3.RouteAction_Cluster{
						Cluster: upstream.Name,
					},
					HostRewriteSpecifier: &envoy_route_v3.RouteAction_HostRewriteLiteral{
						HostRewriteLiteral: upstream.Repository.Host,
					},
					PrefixRewrite: upstream.Route.PathPrefixRewriteValue,
				},
			},
		}

		vhosts.Routes = append(vhosts.Routes, route)
	}

	routeConfig.VirtualHosts = append(routeConfig.VirtualHosts, vhosts)
	data, err := proto.Marshal(http_connection_manager)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize http connection manager: %w", err)
	}

	filterChain := &envoy_listener_v3.FilterChain{
		Filters: make([]*envoy_listener_v3.Filter, 0),
	}

	filter := &envoy_listener_v3.Filter{
		Name: "envoy.filters.network.http_connection_manager",
		ConfigType: &envoy_listener_v3.Filter_TypedConfig{
			TypedConfig: &anypb.Any{
				TypeUrl: "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager",
				Value:   data,
			},
		},
	}

	filterChain.Filters = append(filterChain.Filters, filter)
	listener.FilterChains = append(listener.FilterChains, filterChain)

	return listener, nil
}

func envoyGenerateStaticClusters(gateway *gen.GatewayConfiguration) ([]*envoy_cluster_v3.Cluster, error) {
	clusters := make([]*envoy_cluster_v3.Cluster, 0)

	for _, upstream := range gateway.Upstreams {
		cluster := &envoy_cluster_v3.Cluster{
			Name:            upstream.Name,
			DnsLookupFamily: envoy_cluster_v3.Cluster_V4_ONLY,
			LoadAssignment: &envoy_endpoint_v3.ClusterLoadAssignment{
				ClusterName: upstream.Name,
				Endpoints:   make([]*envoy_endpoint_v3.LocalityLbEndpoints, 0),
			},
			TransportSocket: &envoy_core_v3.TransportSocket{},
		}

		endpoint := &envoy_endpoint_v3.LocalityLbEndpoints{
			LbEndpoints: make([]*envoy_endpoint_v3.LbEndpoint, 0),
		}

		portValue, err := strconv.Atoi(upstream.Repository.Port)
		if err != nil {
			return clusters, fmt.Errorf("failed to convert port from string to uint: %w", err)
		}

		lb_endpoint := &envoy_endpoint_v3.LbEndpoint{
			HostIdentifier: &envoy_endpoint_v3.LbEndpoint_Endpoint{
				Endpoint: &envoy_endpoint_v3.Endpoint{
					Address: &envoy_core_v3.Address{
						Address: &envoy_core_v3.Address_SocketAddress{
							SocketAddress: &envoy_core_v3.SocketAddress{
								Address:       upstream.Repository.Host,
								PortSpecifier: &envoy_core_v3.SocketAddress_PortValue{PortValue: uint32(portValue)},
							},
						},
					},
				},
			},
		}

		endpoint.LbEndpoints = append(endpoint.LbEndpoints, lb_endpoint)
		cluster.LoadAssignment.Endpoints = append(cluster.LoadAssignment.Endpoints, endpoint)

		upstreamTlsTypeConfig := &envoy_extension_tls_v3.UpstreamTlsContext{
			Sni: upstream.Repository.Sni,
		}

		upstreamTlsTypeConfigBinary, err := proto.Marshal(upstreamTlsTypeConfig)
		if err != nil {
			return clusters, fmt.Errorf("failed to serialize to binary: %w", err)
		}

		if upstream.Repository.Tls {
			cluster.TransportSocket = &envoy_core_v3.TransportSocket{
				Name: "envoy.transport_sockets.tls",
				ConfigType: &envoy_core_v3.TransportSocket_TypedConfig{
					TypedConfig: &anypb.Any{
						TypeUrl: "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext",
						Value:   upstreamTlsTypeConfigBinary,
					},
				},
			}
		}

		clusters = append(clusters, cluster)
	}

	pdpCluster, err := envoyGetInternalServiceCluster(pdpSvcName, pdpHostName, PdpPort)
	if err != nil {
		return clusters, fmt.Errorf("failed to generate pdp cluster: %w", err)
	}

	clusters = append(clusters, pdpCluster)

	tapCluster, err := envoyGetInternalServiceCluster(tapSvcName, tapHostName, tapPort)
	if err != nil {
		return clusters, fmt.Errorf("failed to generate tap cluster: %w", err)
	}

	clusters = append(clusters, tapCluster)

	return clusters, nil
}

func envoyGetInternalServiceCluster(svcName, hostName, port string) (*envoy_cluster_v3.Cluster, error) {
	cluster := &envoy_cluster_v3.Cluster{
		Name: svcName,
		LoadAssignment: &envoy_endpoint_v3.ClusterLoadAssignment{
			ClusterName: svcName,
			Endpoints:   make([]*envoy_endpoint_v3.LocalityLbEndpoints, 0),
		},
		TypedExtensionProtocolOptions: map[string]*anypb.Any{},
	}

	explicitHttpConfig := &envoy_extension_http_v3.HttpProtocolOptions{
		UpstreamProtocolOptions: &envoy_extension_http_v3.HttpProtocolOptions_ExplicitHttpConfig_{
			ExplicitHttpConfig: &envoy_extension_http_v3.HttpProtocolOptions_ExplicitHttpConfig{
				ProtocolConfig: &envoy_extension_http_v3.HttpProtocolOptions_ExplicitHttpConfig_Http2ProtocolOptions{
					Http2ProtocolOptions: &envoy_core_v3.Http2ProtocolOptions{},
				},
			},
		},
	}

	explicitHttpConfigData, err := proto.Marshal(explicitHttpConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize proto message: %w", err)
	}

	cluster.TypedExtensionProtocolOptions["envoy.extensions.upstreams.http.v3.HttpProtocolOptions"] = &anypb.Any{
		TypeUrl: "type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions",
		Value:   explicitHttpConfigData,
	}

	endpoint := &envoy_endpoint_v3.LocalityLbEndpoints{
		LbEndpoints: make([]*envoy_endpoint_v3.LbEndpoint, 0),
	}

	portValue, err := strconv.Atoi(port)
	if err != nil {
		return cluster, fmt.Errorf("failed to convert port from string to uint: %w", err)
	}

	lb_endpoint := &envoy_endpoint_v3.LbEndpoint{
		HostIdentifier: &envoy_endpoint_v3.LbEndpoint_Endpoint{
			Endpoint: &envoy_endpoint_v3.Endpoint{
				Address: &envoy_core_v3.Address{
					Address: &envoy_core_v3.Address_SocketAddress{
						SocketAddress: &envoy_core_v3.SocketAddress{
							Address:       hostName,
							PortSpecifier: &envoy_core_v3.SocketAddress_PortValue{PortValue: uint32(portValue)},
						},
					},
				},
			},
		},
	}

	endpoint.LbEndpoints = append(endpoint.LbEndpoints, lb_endpoint)
	cluster.LoadAssignment.Endpoints = append(cluster.LoadAssignment.Endpoints, endpoint)

	return cluster, nil
}

func envoyNodeId(gid string) string {
	return fmt.Sprintf("%s--%s", gid, utils.NewUniqueId())
}
