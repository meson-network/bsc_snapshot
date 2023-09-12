package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/meson-network/bsc-data-file-utils/cmd_download"
	"github.com/meson-network/bsc-data-file-utils/cmd_endpoint"
	"github.com/meson-network/bsc-data-file-utils/cmd_split"
	"github.com/meson-network/bsc-data-file-utils/cmd_upload"
)

const CMD_NAME_SPLIT = "split"
const CMD_NAME_DOWNLOAD = "download"
const CMD_NAME_UPLOAD = "upload"
const CMD_NAME_ENDPOINT = "endpoint"

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
				Subcommands: []*cli.Command{
					{
						Name:  "r2",
						Flags: cmd_upload.GetFlags(),
						Usage: "upload to cloudflare R2 storage",
						Action: func(clictx *cli.Context) error {
							cmd_upload.Upload_r2(clictx)
							return nil
						},
					},
				},
			},
			{
				Name:  CMD_NAME_ENDPOINT,
				Usage: "set endpoint",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Flags: cmd_endpoint.GetFlags(),
						Usage: "add new endpoints",
						Action: func(clictx *cli.Context) error {
							cmd_endpoint.AddEndpoint(clictx)
							return nil
						},
					},
					{
						Name:  "remove",
						Flags: cmd_endpoint.GetFlags(),
						Usage: "remove endpoints",
						Action: func(clictx *cli.Context) error {
							cmd_endpoint.RemoveEndpoint(clictx)
							return nil
						},
					},
					{
						Name:  "set",
						Flags: cmd_endpoint.GetFlags(),
						Usage: "reset endpoints",
						Action: func(clictx *cli.Context) error {
							cmd_endpoint.SetEndpoint(clictx)
							return nil
						},
					},
					{
						Name:  "clear",
						Flags: cmd_endpoint.GetFlags(),
						Usage: "remove all exist endpoints",
						Action: func(clictx *cli.Context) error {
							cmd_endpoint.ClearEndpoint(clictx)
							return nil
						},
					},
					{
						Name:  "print",
						Flags: cmd_endpoint.GetFlags(),
						Usage: "remove all exist endpoints",
						Action: func(clictx *cli.Context) error {
							cmd_endpoint.PrintEndpoint(clictx)
							return nil
						},
					},
				},
			},
		},
	}
}
