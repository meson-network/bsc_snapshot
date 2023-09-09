package cmd_db

import (
	"github.com/meson-network/bsc-data-file-utils/basic/config"
	"github.com/meson-network/bsc-data-file-utils/plugin/g_counter_plugin"
	"github.com/meson-network/bsc-data-file-utils/plugin/sqldb_plugin"
	"github.com/meson-network/bsc-data-file-utils/src/common/dbkv"
	"github.com/meson-network/bsc-data-file-utils/src/user_mgr"
)

func Migrate() {
	StartDBComponent()

	//upgrade table
	sqldb_plugin.GetInstance().AutoMigrate(&user_mgr.UserModel{}, &dbkv.DBKVModel{})

	//try to init if enabled as database table may need to be created
	toml_conf := config.Get_config().Toml_config
	if toml_conf.General_counter.Enable {
		g_counter_plugin.DBMigrate(sqldb_plugin.GetInstance())
	}

}
