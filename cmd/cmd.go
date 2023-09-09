package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/meson-network/bsc-data-file-utils/cmd_default/http/api"
	"github.com/meson-network/bsc-data-file-utils/cmd_split"
	"github.com/meson-network/bsc-data-file-utils/cmd_upload"
)

const CMD_NAME_SPLIT = "split"
const CMD_NAME_DOWNLOAD = "download"
const CMD_NAME_UPLOAD = "upload"
const CMD_NAME_MERGE = "merge"

// const CMD_NAME_DEFAULT = "default"
// const CMD_NAME_GEN_API = "gen_api"
// const CMD_NAME_LOG = "log"
// const CMD_NAME_DB = "db"
// const CMD_NAME_GEOIP = "geoip"
// const CMD_NAME_CONFIG = "config"

// //////config to do cmd ///////////
func ConfigCmd() *cli.App {

	// real_args := config.ConfigBasic("default")

	// var defaultAction = func(clictx *cli.Context) error {
	// 	cmd_default.StartDefault()
	// 	return nil
	// }

	// if len(real_args) > 1 {
	// 	defaultAction = nil
	// }

	return &cli.App{
		// Action: defaultAction, //only run if no sub command

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
					return cmd_split.Split(clictx)
				},
			},
			{
				Name:  CMD_NAME_DOWNLOAD,
				Usage: "multithread download and merge files",
				Flags: cmd_split.GetFlags(),
				Action: func(clictx *cli.Context) error {
					api.GenApiDocs()
					return nil
				},
			},
			{
				Name:  CMD_NAME_UPLOAD,
				Usage: "upload files",
				Flags: cmd_upload.GetFlags(),
				Action: func(clictx *cli.Context) error {
					return cmd_upload.Upload_r2(clictx)
				},
			},

			// {
			// 	Name:  CMD_NAME_GEN_API,
			// 	Usage: "api command",
			// 	Action: func(clictx *cli.Context) error {
			// 		api.GenApiDocs()
			// 		return nil
			// 	},
			// },
			// {
			// 	Name:  CMD_NAME_LOG,
			// 	Usage: "print all logs",
			// 	Flags: cmd_log.GetFlags(),
			// 	Action: func(clictx *cli.Context) error {
			// 		num := clictx.Int64("num")
			// 		onlyerr := clictx.Bool("only_err")
			// 		cmd_log.StartLog(onlyerr, num)
			// 		return nil
			// 	},
			// },
			// {
			// 	Name:  CMD_NAME_DB,
			// 	Usage: "db command",
			// 	Subcommands: []*cli.Command{
			// 		{
			// 			Name:  "init",
			// 			Usage: "initialize db data",
			// 			Action: func(clictx *cli.Context) error {
			// 				fmt.Println("======== start of db data initialization ========")
			// 				cmd_db.Initialize()
			// 				fmt.Println("======== end  of  db data initialization ========")
			// 				return nil
			// 			},
			// 		},
			// 		{
			// 			Name:  "migrate",
			// 			Usage: "migrate db tables",
			// 			Action: func(clictx *cli.Context) error {
			// 				fmt.Println("======== start of db table migration ========")
			// 				cmd_db.Migrate()
			// 				fmt.Println("======== end of db table migration ========")
			// 				return nil
			// 			},
			// 		},
			// 	},
			// },
			// {
			// 	Name:  CMD_NAME_GEOIP,
			// 	Usage: "geoip command",
			// 	Subcommands: []*cli.Command{
			// 		{
			// 			Name:  "download",
			// 			Usage: "download geoip data",
			// 			Action: func(clictx *cli.Context) error {
			// 				fmt.Println("======== start of geoip data download ========")
			// 				cmd_geoip.Download()
			// 				fmt.Println("======== end  of  geoip data download ========")
			// 				return nil
			// 			},
			// 		},
			// 	},
			// },
			// {
			// 	Name:  CMD_NAME_CONFIG,
			// 	Usage: "config command",
			// 	Subcommands: []*cli.Command{
			// 		//show config
			// 		{
			// 			Name:  "get_toml",
			// 			Usage: "show config as toml string",
			// 			Action: func(clictx *cli.Context) error {
			// 				configs, _ := config.Get_config().MergedConfig2TomlStr()
			// 				fmt.Println(configs)
			// 				return nil
			// 			},
			// 		},
			// 		//show config
			// 		{
			// 			Name:  "get_json",
			// 			Usage: "show config as json string",
			// 			Action: func(clictx *cli.Context) error {
			// 				configs, _ := config.Get_config().MergedConfig2JsonStr()
			// 				fmt.Println(configs)
			// 				return nil
			// 			},
			// 		},
			// 		//set config
			// 		{
			// 			Name:  "set",
			// 			Usage: "set config",
			// 			Flags: append(cmd_conf.Cli_get_flags(), &cli.StringFlag{Name: "config", Required: false}),
			// 			Action: func(clictx *cli.Context) error {
			// 				return cmd_conf.Cli_set_config(clictx)
			// 			},
			// 		},
			// 	},
			// },
		},
	}
}
