package component

import (
	"errors"
	"strings"

	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/basic/config"
	"github.com/meson-network/bsc-data-file-utils/plugin/sqlite_plugin"
)

func InitSqlite(toml_conf *config.TomlConfig) error {

	if toml_conf.Sqlite.Enable {

		var sqlite_conf sqlite_plugin.Config

		if !strings.Contains(toml_conf.Sqlite.Path, "file::memory") {
			//file check
			sqlite_abs_path, sqlite_abs_path_exist, _ := basic.PathExist(toml_conf.Sqlite.Path)
			if !sqlite_abs_path_exist {
				return errors.New(toml_conf.Sqlite.Path + " :sqlite.path not exist , please reset your sqlite.path :" + toml_conf.Sqlite.Path)
			}
			sqlite_conf = sqlite_plugin.Config{
				Sqlite_abs_path: sqlite_abs_path,
			}
		} else {
			sqlite_conf = sqlite_plugin.Config{
				Sqlite_abs_path: toml_conf.Sqlite.Path,
			}
		}

		basic.Logger.Infoln("Init sqlite plugin with config:", sqlite_conf)
		if err := sqlite_plugin.Init(&sqlite_conf, basic.Logger); err == nil {
			basic.Logger.Infoln("### InitSqlite success")
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}
