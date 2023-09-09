package component

import (
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/basic/config"
	"github.com/meson-network/bsc-data-file-utils/plugin/sqldb_plugin"
)

func InitDB(toml_conf *config.TomlConfig) error {

	if toml_conf.Db.Enable {
		db_conf := sqldb_plugin.Config{
			Host:     toml_conf.Db.Host,
			Port:     toml_conf.Db.Port,
			DbName:   toml_conf.Db.Name,
			UserName: toml_conf.Db.Username,
			Password: toml_conf.Db.Password,
		}
		basic.Logger.Infoln("Init db plugin with config:", db_conf)
		if err := sqldb_plugin.Init(&db_conf, basic.Logger); err == nil {
			basic.Logger.Infoln("### InitDB success")
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}
