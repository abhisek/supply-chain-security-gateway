package openssf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	backoff := heimdall.NewExponentialBackoff(config.Timeout/time.Second,
		config.Timeout, 0.5, time.Millisecond*100)

	client := hystrix.NewClient(
		hystrix.WithHTTPTimeout(config.Timeout),
		hystrix.WithCommandName("osv_api_get_request"),
		hystrix.WithMaxConcurrentRequests(config.MaxConcurrentRequest),
		hystrix.WithRetryCount(config.Retry),
		hystrix.WithRetrier(heimdall.NewRetrier(backoff)),
	)

	return &OsvServiceAdapter{config: config, client: client}
}

func (svc *OsvServiceAdapter) QueryPackage(ecosystem, name string) (V1VulnerabilityList, error) {
	rQuery := &V1Query{}
	rQuery.Package.Ecosystem = &ecosystem
	rQuery.Package.Name = &name

	body, err := json.Marshal(rQuery)
	if err != nil {
		return V1VulnerabilityList{}, err
	}

	resp, err := svc.client.Post(OsvQueryEndpoint, bytes.NewReader(body), http.Header{})
	if err != nil {
		return V1VulnerabilityList{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return V1VulnerabilityList{}, fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}

	var vulnList V1VulnerabilityList
	err = json.NewDecoder(resp.Body).Decode(&vulnList)
	if err != nil {
		return V1VulnerabilityList{}, fmt.Errorf("failed to decoded to vulnerability list: %v", err)
	}

	return vulnList, nil
}
