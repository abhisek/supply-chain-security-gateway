package utils

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// Serialize an interface using JSON or return error string
func Introspect(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	} else {
		return string(bytes)
	}
}

func MapStruct[T any](source interface{}, dest *T) error {
	return mapstructure.Decode(source, dest)
}
