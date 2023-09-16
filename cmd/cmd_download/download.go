package cmd_download

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/meson-network/bsc_snapshot/basic"
	"github.com/meson-network/bsc_snapshot/basic/color"
	"github.com/meson-network/bsc_snapshot/src/download"
)

func Download(clictx *cli.Context) error {
	fmt.Println(color.Green(basic.Logo))

	jsonConfigAddress, thread, retryNum, noResume := ReadParam(clictx)
	return download.Download(jsonConfigAddress, thread, retryNum, noResume)
}
