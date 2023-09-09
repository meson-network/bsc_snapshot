package cmd_conf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/meson-network/bsc-data-file-utils/basic/color"
	"github.com/meson-network/bsc-data-file-utils/basic/config"

	"github.com/urfave/cli/v2"
)

func Cli_get_flags() []cli.Flag {
	allflags := []cli.Flag{}
	allflags = append(allflags, &cli.StringFlag{Name: "build.mode", Required: false})
	allflags = append(allflags, &cli.StringFlag{Name: "log.level", Required: false})
	allflags = append(allflags, &cli.BoolFlag{Name: "http.enable", Required: false})
	allflags = append(allflags, &cli.IntFlag{Name: "http.port", Required: false})
	allflags = append(allflags, &cli.BoolFlag{Name: "https.enable", Required: false})
	allflags = append(allflags, &cli.BoolFlag{Name: "db.enable", Required: false})
	allflags = append(allflags, &cli.BoolFlag{Name: "redis.enable", Required: false})

	return allflags
}

var log_level_map = map[string]struct{}{
	"TRAC":  {},
	"TRACE": {},
	"DEBU":  {},
	"DEBUG": {},
	"INFO":  {},
	"WARN":  {},
	"ERRO":  {},
	"ERROR": {},
	"FATA":  {},
	"FATAL": {},
	"PANI":  {},
	"PANIC": {},
}

const BUILD_MODE_RELEASE = "release"
const BUILD_MODE_DEBUG = "debug"

func Cli_set_config(clictx *cli.Context) error {
	config := config.Get_config()

	if clictx.IsSet("build.mode") {
		build_m := clictx.String("build.mode")
		build_m = strings.TrimSpace(strings.ToLower(build_m))
		if build_m != BUILD_MODE_RELEASE && build_m != BUILD_MODE_DEBUG {
			return errors.New("build.mode error, allowed: '" + BUILD_MODE_RELEASE + "' or '" + BUILD_MODE_DEBUG + "'")
		}
		config.SetUserTomlConfig("build.mode", build_m)
	}

	if clictx.IsSet("log.level") {
		log_level := strings.ToUpper(clictx.String("log.level"))
		_, exist := log_level_map[log_level]
		if !exist {
			return errors.New("log level error")
		}

		config.SetUserTomlConfig("log.level", log_level)
	}

	if clictx.IsSet("http.enable") {
		http_enable := clictx.Bool("http.enable")
		config.SetUserTomlConfig("http.enable", http_enable)
	}

	if clictx.IsSet("http.port") {
		http_port := clictx.Int64("http.port")
		config.SetUserTomlConfig("http.port", http_port)
	}

	if clictx.IsSet("https.enable") {
		https_enable := clictx.Bool("https.enable")
		config.SetUserTomlConfig("https.enable", https_enable)
	}

	if clictx.IsSet("db.enable") {
		db_enable := clictx.Bool("db.enable")
		config.SetUserTomlConfig("db.enable", db_enable)
	}

	if clictx.IsSet("redis.enable") {
		redis_enable := clictx.Bool("redis.enable")
		config.SetUserTomlConfig("redis.enable", redis_enable)
	}

	err := config.Save_user_config()
	if err != nil {
		fmt.Println(color.Red("save user config error:", err))
		return err
	} else {
		fmt.Println(color.Green("new config generated"))
		fmt.Println(color.Green("restart for the new configuration to take effect"))
	}
	return nil
}
