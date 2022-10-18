package openssf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/logger"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/hystrix"
)

const (
	OsvQueryEndpoint = "https://api.osv.dev/v1/query"
)

type OsvServiceAdapterConfig struct {
	Timeout              time.Duration
	Retry                int
	MaxConcurrentRequest int
}

func DefaultServiceAdapterConfig() OsvServiceAdapterConfig {
	return OsvServiceAdapterConfig{
		Timeout:              5 * time.Second,
		Retry:                3,
		MaxConcurrentRequest: 5,
	}
}

type OsvServiceAdapter struct {
	config OsvServiceAdapterConfig
	client *hystrix.Client
}

func NewOsvServiceAdapter(config OsvServiceAdapterConfig) *OsvServiceAdapter {
	backoff := heimdall.NewConstantBackoff(2*time.Second, 100*time.Millisecond)

	client := hystrix.NewClient(
		hystrix.WithHTTPTimeout(config.Timeout),
		hystrix.WithCommandName("osv_api_get_request"),
		hystrix.WithMaxConcurrentRequests(config.MaxConcurrentRequest),
		hystrix.WithRetryCount(config.Retry),
		hystrix.WithRetrier(heimdall.NewRetrier(backoff)),
	)

	return &OsvServiceAdapter{config: config, client: client}
}

func (svc *OsvServiceAdapter) QueryPackage(ecosystem, name, version string) (V1VulnerabilityList, error) {
	rQuery := &V1Query{
		Package: &struct {
			OsvPackage "yaml:\",inline\""
		}{},
	}

	rQuery.Version = &version
	rQuery.Package.Ecosystem = &ecosystem
	rQuery.Package.Name = &name

	logger.Debugf("Querying OSV with: ecosystem:%s name:%s version:%s",
		ecosystem, name, version)

	body, err := json.Marshal(rQuery)
	if err != nil {
		return V1VulnerabilityList{}, err
	}

	resp, err := svc.client.Post(OsvQueryEndpoint, bytes.NewReader(body), http.Header{})
	if err != nil {
		return V1VulnerabilityList{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return V1VulnerabilityList{}, fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}

	var vulnList V1VulnerabilityList
	err = json.NewDecoder(resp.Body).Decode(&vulnList)
	if err != nil {
		return V1VulnerabilityList{}, fmt.Errorf("failed to decoded to vulnerability list: %w", err)
	}

	// No result found
	if vulnList.Vulns == nil {
		return V1VulnerabilityList{}, fmt.Errorf("empty vulnerability list from OSV")
	}

	return vulnList, nil
}
