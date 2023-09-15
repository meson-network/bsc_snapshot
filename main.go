package main

import (
	"os"

	"github.com/meson-network/bsc_snapshot/cmd"
)

func main() {

	//config app to run
	errRun := cmd.ConfigCmd().Run(os.Args)
	if errRun != nil {
		os.Exit(1)
	}
}
