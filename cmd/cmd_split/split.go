package cmd_split

import (
	"github.com/urfave/cli/v2"

	"github.com/meson-network/bsc_snapshot/src/split"
)

func Split(clictx *cli.Context) error {
	originFilePath, destDir, sizeStr, thread := ReadParam(clictx)

	return split.Split(originFilePath, destDir, sizeStr, thread)
}
