package cmd_download

import "github.com/urfave/cli/v2"

func GetFlags() (allflags []cli.Flag) {
	return []cli.Flag{
		&cli.StringFlag{Name: "file_config", Required: true},
		&cli.StringFlag{Name: "thread", Required: false},
		&cli.IntFlag{Name: "retry_times", Required: false},
		&cli.BoolFlag{Name: "no_resume", Required: false},
	}
}
