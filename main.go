package main

import (
	"os"

	"github.com/meson-network/bsc-data-file-utils/cmd"
)

func main() {

	//config app to run
	errRun := cmd.ConfigCmd().Run(os.Args)
	if errRun != nil {
		panic(errRun)
	}
}
