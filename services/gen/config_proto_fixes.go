package gen

import (
	"encoding/json"
	"fmt"
)

func (s *GatewayUpstreamType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	err := unmarshal(&v)
	if err != nil {
		return err
	}

	if iv, ok := GatewayUpstreamType_value[v]; !ok {
		return fmt.Errorf("unknown upstream type: %s", v)
	} else {
		*s = GatewayUpstreamType(iv)
	}

	return nil
}

func (s *GatewayUpstreamType) MarshalYAML() (interface{}, error) {
	return GatewayUpstreamType_name[int32(*s)], nil
}

func (s *GatewayUpstreamType) MarshalJSON() ([]byte, error) {
	return json.Marshal(GatewayUpstreamType_name[int32(*s)])
}

func (s *GatewayAuthenticationType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	err := unmarshal(&v)
	if err != nil {
		return err
	}

	if iv, ok := GatewayAuthenticationType_value[v]; !ok {
		return fmt.Errorf("unknown authentication type: %s", v)
	} else {
		*s = GatewayAuthenticationType(iv)
	}

	return nil
}

func (s *GatewayAuthenticationType) MarshalYAML() (interface{}, error) {
	return GatewayAuthenticationType_name[int32(*s)], nil
}

func (s *GatewayAuthenticationType) MarshalJSON() ([]byte, error) {
	return json.Marshal(GatewayAuthenticationType_name[int32(*s)])
}

func (s *GatewaySecretSource) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	err := unmarshal(&v)
	if err != nil {
		return err
	}

	if iv, ok := GatewaySecretSource_value[v]; !ok {
		return fmt.Errorf("unknown secret source: %s", v)
	} else {
		*s = GatewaySecretSource(iv)
	}

	return nil
}

func (s *GatewaySecretSource) MarshalYAML() (interface{}, error) {
	return GatewaySecretSource_name[int32(*s)], nil
}

func (s *GatewaySecretSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(GatewaySecretSource_name[int32(*s)])
}
