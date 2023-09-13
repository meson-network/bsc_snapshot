package file_config

import (
	"errors"
	"strings"
)

func FormatEndpoints(endpoints []string) ([]string, error) {
	for i := range endpoints {
		endpoints[i] = strings.TrimSpace(endpoints[i])
		if !strings.HasPrefix(endpoints[i], "http") {
			return nil, errors.New("endpoint format error")
		}
		for strings.HasSuffix(endpoints[i], "/") {
			endpoints[i] = strings.TrimSuffix(endpoints[i], "/")
		}
	}

	return endpoints, nil
}
