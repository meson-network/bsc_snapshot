package component

import (
	"errors"

	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/basic/config"
	"github.com/meson-network/bsc-data-file-utils/plugin/echo_plugin"
)

func Init_http_echo_server(toml_conf *config.TomlConfig) error {

	if toml_conf.Http.Enable {
		if err := echo_plugin.Init_("http", echo_plugin.Config{Port: toml_conf.Http.Port, Tls: false, Crt_path: "", Key_path: ""},
			func(panic_err interface{}) {
				basic.Logger.Errorln(panic_err)
			}, basic.Logger); err == nil {
			basic.Logger.Infoln("### Init_http_echo_server success")
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}

func Init_https_echo_server(toml_conf *config.TomlConfig) error {

	if toml_conf.Https.Enable {

		crt_abs_path, crt_path_exist, _ := basic.PathExist(toml_conf.Https.Crt_path)
		if !crt_path_exist {
			return errors.New("https crt file path error:" + toml_conf.Https.Crt_path)
		}

		key_abs_path, key_path_exist, _ := basic.PathExist(toml_conf.Https.Key_path)
		if !key_path_exist {
			return errors.New("https key file path error:" + toml_conf.Https.Key_path)
		}

		if err := echo_plugin.Init_("https", echo_plugin.Config{Port: toml_conf.Https.Port, Tls: true, Crt_path: crt_abs_path, Key_path: key_abs_path},
			func(panic_err interface{}) {
				basic.Logger.Errorln(panic_err)
			}, basic.Logger); err == nil {
			basic.Logger.Infoln("### Init_https_echo_server success")
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}

func InitEchoServer(toml_conf *config.TomlConfig) error {
	http_err := Init_http_echo_server(toml_conf)
	if http_err != nil {
		return http_err
	}
	https_err := Init_https_echo_server(toml_conf)
	if https_err != nil {
		return https_err
	}

	return nil
}
