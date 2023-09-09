package cmd_split

import "github.com/urfave/cli/v2"

func GetFlags() (allflags []cli.Flag) {
	return []cli.Flag{
		&cli.StringFlag{Name: "file", Required: true},
		&cli.StringFlag{Name: "dest", Required: false},
		&cli.StringFlag{Name: "size", Required: true},
	}
}
