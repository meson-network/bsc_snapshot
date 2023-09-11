package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/meson-network/bsc-data-file-utils/cmd_download"
	"github.com/meson-network/bsc-data-file-utils/cmd_split"
	"github.com/meson-network/bsc-data-file-utils/cmd_upload"
)

const CMD_NAME_SPLIT = "split"
const CMD_NAME_DOWNLOAD = "download"
const CMD_NAME_UPLOAD = "upload"
const CMD_NAME_MERGE = "merge"

// //////config to do cmd ///////////
func ConfigCmd() *cli.App {

	return &cli.App{

		//run if sub command not correct
		CommandNotFound: func(context *cli.Context, s string) {
			fmt.Println("command not find, use -h or --help show help")
		},

		Commands: []*cli.Command{
			{
				Name:  CMD_NAME_SPLIT,
				Usage: "split data file to small files",
				Flags: cmd_split.GetFlags(),
				Action: func(clictx *cli.Context) error {
					cmd_split.Split(clictx)
					return nil
				},
			},
			{
				Name:  CMD_NAME_DOWNLOAD,
				Usage: "multithread download and merge files",
				Flags: cmd_download.GetFlags(),
				Action: func(clictx *cli.Context) error {
					cmd_download.Download(clictx)
					return nil
				},
			},
			{
				Name:  CMD_NAME_UPLOAD,
				Usage: "upload files",
				Flags: cmd_upload.GetFlags(),
				Action: func(clictx *cli.Context) error {
					cmd_upload.Upload_r2(clictx)
					return nil
				},
			},

		},
	}
}
