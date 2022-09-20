package main

import (
	"encoding/json"
	"os"

	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/logger"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
	"github.com/ghodss/yaml"
)

const (
	ingressGatewayBasicAuthenticatorName = "default-basic-auth"
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

	s.printConfig(gateway)

	return nil
}

// We serialize to JSON first because proto generated classes has JSON
// key name annotations
func (s *sampleConfigGenerator) printConfig(gateway *config_api.GatewayConfiguration) {
	data, err := json.Marshal(gateway)
	if err != nil {
		logger.Errorf("Failed to JSON serialize gateway config: %v", err)
		return
	}

	yamlData, err := yaml.JSONToYAML(data)
	if err != nil {
		logger.Errorf("Failed to convert YAML: %v", err)
		return
	}

	os.Stdout.Write(yamlData)
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
			BasicAuth: &config_api.GatewayAuthenticatorBasicAuth{
				Path: "data/default-basic-auth",
			},
		},
	}
}

func (s *sampleConfigGenerator) addDefaultUpstreams(gateway *config_api.GatewayConfiguration) {
	gateway.Upstreams = []*config_api.GatewayUpstream{}

	gateway.Upstreams = append(gateway.Upstreams,
		s.getUpstream("maven-central", config_api.GatewayUpstreamType_Maven,
			config_api.GatewayUpstreamManagementType_GatewayAdmin, "/maven2",
			"repo.maven.apache.org", "443"))

	gateway.Upstreams = append(gateway.Upstreams, s.getUpstream("gradle-plugins", config_api.GatewayUpstreamType_Maven,
		config_api.GatewayUpstreamManagementType_GatewayAdmin, "/gradle-plugins/m2",
		"plugins.gradle.org", "443"))
}

func (s *sampleConfigGenerator) getUpstream(name string,
	uType config_api.GatewayUpstreamType, mType config_api.GatewayUpstreamManagementType,
	pathPrefix string, host string, port string) *config_api.GatewayUpstream {

	return &config_api.GatewayUpstream{
		Name:           name,
		Type:           uType,
		ManagementType: mType,
		Authentication: &config_api.GatewayAuthenticationProvider{
			Type:     config_api.GatewayAuthenticationType_Basic,
			Provider: ingressGatewayBasicAuthenticatorName,
		},
		Route: &config_api.GatewayUpstreamRoute{
			PathPrefix: pathPrefix,
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
