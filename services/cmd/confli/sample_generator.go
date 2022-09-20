package main

import (
	"encoding/json"
	"os"

	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/logger"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
	"github.com/ghodss/yaml"
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

func (s *sampleConfigGenerator) addDefaultUpstreams(gateway *config_api.GatewayConfiguration) {
	gateway.Upstreams = []*config_api.GatewayUpstream{}

	gateway.Upstreams = append(gateway.Upstreams, &config_api.GatewayUpstream{
		Name:           "mavenCentral",
		Type:           config_api.GatewayUpstreamType_Maven,
		ManagementType: config_api.GatewayUpstreamManagementType_GatewayAdmin,
		Authentication: &config_api.GatewayAuthenticationProvider{},
		Route: &config_api.GatewayUpstreamRoute{
			PathPrefix: "/maven2",
		},
		Repository: &config_api.GatewayUpstreamRepository{
			Host:           "repo.maven.apache.org",
			Port:           "443",
			Tls:            true,
			Sni:            "repo.maven.apache.org",
			Authentication: &config_api.GatewayAuthenticationProvider{},
		},
	})
}
