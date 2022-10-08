package main

import (
	"os"

	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/logger"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
	"github.com/golang/protobuf/jsonpb"
)

const (
	ingressGatewayBasicAuthenticatorName = "default-basic-auth"
	messagingAdapterNameNATS             = "nats"
	messagingAdapterNameKafka            = "kafka"
)

type sampleConfigGenerator struct {
	file string
}

func newSampleConfigGenerator(path string) *sampleConfigGenerator {
	return &sampleConfigGenerator{file: path}
}

func (s *sampleConfigGenerator) generate() error {
	gateway := &config_api.GatewayConfiguration{}

	s.addInfo(gateway)
	s.addDefaultUpstreams(gateway)
	s.addDefaultGatewayAuth(gateway)
	s.addMessaging(gateway)
	s.addServiceConfig(gateway)

	s.printConfig(gateway)

	return nil
}

// We serialize to JSON first because proto generated classes has JSON
// key name annotations
func (s *sampleConfigGenerator) printConfig(gateway *config_api.GatewayConfiguration) {
	m := jsonpb.Marshaler{Indent: "  "}
	data, err := m.MarshalToString(gateway)
	if err != nil {
		logger.Errorf("Failed to JSON serialize gateway config: %v", err)
		return
	}

	os.Stdout.Write([]byte(data))
}

func (s *sampleConfigGenerator) addInfo(gateway *config_api.GatewayConfiguration) {
	gateway.Info = &config_api.GatewayInfo{
		Id:     utils.NewUniqueId(),
		Name:   gatewayName,
		Domain: gatewayDomain,
	}
}

func (s *sampleConfigGenerator) addDefaultGatewayAuth(gateway *config_api.GatewayConfiguration) {
	gateway.Authenticators = map[string]*config_api.GatewayAuthenticator{
		ingressGatewayBasicAuthenticatorName: {
			Type: config_api.GatewayAuthenticationType_Basic,
			Config: &config_api.GatewayAuthenticator_BasicAuth{
				BasicAuth: &config_api.GatewayAuthenticatorBasicAuth{
					Path: "data/default-basic-auth.txt",
				},
			},
		},
	}
}

func (s *sampleConfigGenerator) addDefaultUpstreams(gateway *config_api.GatewayConfiguration) {
	gateway.Upstreams = []*config_api.GatewayUpstream{}

	gateway.Upstreams = append(gateway.Upstreams,
		s.getUpstream("maven-central", config_api.GatewayUpstreamType_Maven,
			config_api.GatewayUpstreamManagementType_GatewayAdmin, "/maven2", "/maven2",
			"repo.maven.apache.org", "443"))

	gateway.Upstreams = append(gateway.Upstreams, s.getUpstream("gradle-plugins", config_api.GatewayUpstreamType_Maven,
		config_api.GatewayUpstreamManagementType_GatewayAdmin, "/gradle-plugins/m2", "/m2",
		"plugins.gradle.org", "443"))
}

func (s *sampleConfigGenerator) addMessaging(gateway *config_api.GatewayConfiguration) {
	gateway.Messaging = map[string]*config_api.MessagingAdapter{
		messagingAdapterNameNATS: {
			Type: config_api.MessagingAdapter_NATS,
			Config: &config_api.MessagingAdapter_Nats{
				Nats: &config_api.MessagingAdapter_NatsAdapterConfig{
					Url: natsUrl,
				},
			},
		},
		messagingAdapterNameKafka: {
			Type: config_api.MessagingAdapter_KAFKA,
			Config: &config_api.MessagingAdapter_Kafka{
				Kafka: &config_api.MessagingAdapter_KafkaAdapterConfig{
					BootstrapServers:  []string{"kafka-host:9092"},
					SchemaRegistryUrl: "http://kafka-host:8081",
				},
			},
		},
	}
}

func (s *sampleConfigGenerator) addServiceConfig(gateway *config_api.GatewayConfiguration) {
	gateway.Services = &config_api.GatewayConfiguration_ServiceConfig{}

	s.addPdpServiceConfig(gateway.Services)
	s.addTapServiceConfig(gateway.Services)
	s.addDcsServiceConfig(gateway.Services)
}

func (s *sampleConfigGenerator) addPdpServiceConfig(config *config_api.GatewayConfiguration_ServiceConfig) {
	config.Pdp = &config_api.PdpServiceConfig{
		MonitorMode: true,
		PublisherConfig: &config_api.PdpServiceConfig_PublisherConfig{
			MessagingAdapterName: messagingAdapterNameNATS,
			TopicNames: &config_api.PdpServiceConfig_PublisherConfig_TopicNames{
				PolicyAudit: "gateway.pdp.audits",
			},
		},
		PdsClient: &config_api.PdsClientConfig{
			Type: config_api.PdsClientType_LOCAL,
			Config: &config_api.PdsClientConfig_Common{
				Common: &config_api.PdsClientCommonConfig{
					Host: "pds",
					Port: 9002,
					Mtls: true,
				},
			},
		},
	}
}

func (s *sampleConfigGenerator) addTapServiceConfig(config *config_api.GatewayConfiguration_ServiceConfig) {
	config.Tap = &config_api.TapServiceConfig{
		PublisherConfig: &config_api.TapServiceConfig_PublisherConfig{
			MessagingAdapterName: messagingAdapterNameNATS,
			TopicNames: &config_api.TapServiceConfig_PublisherConfig_TopicNames{
				UpstreamRequest:  "gateway.tap.upstream_req",
				UpstreamResponse: "gateway.tap.upstream_res",
			},
		},
	}
}

func (s *sampleConfigGenerator) addDcsServiceConfig(config *config_api.GatewayConfiguration_ServiceConfig) {
	config.Dcs = &config_api.DcsServiceConfig{
		Active:               true,
		MessagingAdapterName: messagingAdapterNameNATS,
	}
}

func (s *sampleConfigGenerator) getUpstream(name string,
	uType config_api.GatewayUpstreamType, mType config_api.GatewayUpstreamManagementType,
	pathPrefix string, pathRewrite string, host string, port string) *config_api.GatewayUpstream {

	return &config_api.GatewayUpstream{
		Name:           name,
		Type:           uType,
		ManagementType: mType,
		Authentication: &config_api.GatewayAuthenticationProvider{
			Type:     config_api.GatewayAuthenticationType_Basic,
			Provider: ingressGatewayBasicAuthenticatorName,
		},
		Route: &config_api.GatewayUpstreamRoute{
			PathPrefix:             pathPrefix,
			HostRewriteValue:       host,
			PathPrefixRewriteValue: pathRewrite,
		},
		Repository: &config_api.GatewayUpstreamRepository{
			Host:           host,
			Port:           port,
			Tls:            true,
			Sni:            host,
			Authentication: &config_api.GatewayAuthenticationProvider{},
		},
	}
}
