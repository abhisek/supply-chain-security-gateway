package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Implement poor man's feature flag based on some data source
// For now, env + convention is a good data source

const (
	ffEnvPrefix = "FF"
)

// Takes a string such as "app_dcs" and determine if
// feature is explicitly disabled by looking up FF_APP_DCS_DISABLED=true
func IsFeatureDisabled(key string) bool {
	key = fmt.Sprintf("%s_%s_disabled", ffEnvPrefix, key)
	key = strings.ToUpper(key)

	if v, b := os.LookupEnv(key); b {
		ret, err := strconv.ParseBool(v)
		if err == nil {
			return ret
		}
	}

	return false
}
