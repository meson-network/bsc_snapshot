package config

import (
	"os"
	"strings"

	ilog "github.com/coreservice-io/log"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/plugin/logger_plugin"
)

// 1.setup the config toml file
// 2.setup the basic working directory
// 3.setup the basic logger
// 4.return the real args
func ConfigBasic(toml_target string) []string {
	//////////init config/////////////
	real_args := []string{}
	for _, arg := range os.Args {
		arg_lower := strings.ToLower(arg)
		if strings.HasPrefix(arg_lower, "-conf=") || strings.HasPrefix(arg_lower, "--conf=") {
			toml_target = strings.TrimPrefix(arg_lower, "--conf=")
			toml_target = strings.TrimPrefix(toml_target, "-conf=")
			continue
		}
		real_args = append(real_args, arg)
	}

	os.Args = real_args
	conf_err := Init_config(toml_target)
	if conf_err != nil {
		panic(conf_err)
	}

	/////set build/////////
	build_mode_err := basic.SetMode(Get_config().Toml_config.Build.Mode)
	if build_mode_err != nil {
		panic(build_mode_err)
	}

	/////set up basic logger ///////
	logs_abs_path := basic.AbsPath("logs")
	if err := logger_plugin.Init(logs_abs_path); err != nil {
		panic(err)
	}

	basic.Logger = logger_plugin.GetInstance()
	loglevel := ilog.ParseLogLevel(Get_config().Toml_config.Log.Level)
	basic.Logger.SetLevel(loglevel)

	return real_args

}
