package cmd_download

import (
	"github.com/urfave/cli/v2"

	"github.com/meson-network/bsc_snapshot/src/download"
)

func Download(clictx *cli.Context) error {

	jsonConfigAddress, thread, retryNum, noResume := ReadParam(clictx)
	return download.Download(jsonConfigAddress, thread, retryNum, noResume)
}
