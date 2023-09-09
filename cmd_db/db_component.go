package cmd_db

import (
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/basic/config"
	"github.com/meson-network/bsc-data-file-utils/component"
	"github.com/meson-network/bsc-data-file-utils/src/common/token_mgr"
)

func StartDBComponent() {

	toml_conf := config.Get_config().Toml_config

	/////////////////////////
	if err := component.InitReference(); err != nil {
		basic.Logger.Fatalln(err)
	}
	/////////////////////////
	if !toml_conf.Db.Enable {
		basic.Logger.Fatalln("db not enabled in config")
	}
	if err := component.InitDB(toml_conf); err != nil {
		basic.Logger.Fatalln(err)
	}
	/////////////////////////
	if !toml_conf.Redis.Enable {
		basic.Logger.Fatalln("redis not enabled in config")
	}
	if err := component.InitRedis(toml_conf); err != nil {
		basic.Logger.Fatalln(err)
	}

	///////////////////////////////
	if err := token_mgr.InitTokenMgr(toml_conf); err != nil {
		basic.Logger.Fatalln(err)
	}
}
