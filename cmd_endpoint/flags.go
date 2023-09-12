package cmd_endpoint

import "github.com/urfave/cli/v2"

func GetFlags() (allflags []cli.Flag) {
	return []cli.Flag{
		&cli.StringFlag{Name: "config_path", Required: true},
		&cli.StringSliceFlag{Name: "endpoint",Aliases: []string{"e"}, Required: false},
	}
}

func readParam(clictx *cli.Context) (string,[]string) {
	//read params
	configPath := clictx.String("config_path")
	endpoints := clictx.StringSlice("endpoint")

	return configPath,endpoints
}