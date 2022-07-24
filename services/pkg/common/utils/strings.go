package utils

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
)

// Serialize an interface using JSON or return error string
func Introspect(v interface{}) string {
	bytes, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	} else {
		return string(bytes)
	}
}

func CleanPath(path string) string {
	return filepath.Clean(path)
}

func MapStruct[T any](source interface{}, dest *T) error {
	return mapstructure.Decode(source, dest)
}

func SafelyGetValue[T any](target *T) T {
	var vi T
	if target != nil {
		vi = *target
	}

	return vi
}

func IsEmptyString(s string) bool {
	return s == ""
}
