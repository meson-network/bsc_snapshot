package cmd_endpoint

import (
	"github.com/urfave/cli/v2"

	"github.com/meson-network/bsc_snapshot/src/endpoint"
)

func AddEndpoint(clictx *cli.Context) error {
	configPath, endpoints := ReadParam(clictx)
	return endpoint.AddEndpoint(configPath, endpoints)
}

func RemoveEndpoint(clictx *cli.Context) error {
	configPath, endpoints := ReadParam(clictx)
	return endpoint.RemoveEndpoint(configPath, endpoints)
}

func SetEndpoint(clictx *cli.Context) error {
	configPath, endpoints := ReadParam(clictx)
	return endpoint.SetEndpoint(configPath, endpoints)
}

func ClearEndpoint(clictx *cli.Context) error {
	configPath, _ := ReadParam(clictx)
	return endpoint.ClearEndpoint(configPath)
}

func PrintEndpoint(clictx *cli.Context) error {
	configPath, _ := ReadParam(clictx)
	return endpoint.PrintEndpoint(configPath)
}
