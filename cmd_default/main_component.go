package cmd_default

import (
	"io"
	"os"

	"github.com/coreservice-io/log"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/basic/config"
	"github.com/meson-network/bsc-data-file-utils/component"
	"github.com/meson-network/bsc-data-file-utils/src/common/token_mgr"
)

func InitComponent_() error {

	// for better ui console disable infoln output to console e.g: logs only to file
	if basic.Logger.GetLevel() <= log.InfoLevel {
		basic.Logger.SetOutput(io.Discard)
		defer basic.Logger.SetOutput(os.Stdout)
	}

	config_toml := config.Get_config().Toml_config

	// ///////////////////////
	if err := component.InitReference(); err != nil {
		return err
	}
	// init token mgr
	if err := token_mgr.InitTokenMgr(config_toml); err != nil {
		return err
	}
	// ///////////////////////
	if err := component.InitDB(config_toml); err != nil {
		return err
	}
	// ///////////////////////
	if err := component.InitRedis(config_toml); err != nil {
		return err
	}
	// ///////////////////////
	if err := component.InitAutoCert(config_toml); err != nil {
		return err
	}
	// ///////////////////////
	if err := component.InitEchoServer(config_toml); err != nil {
		return err
	}
	// ///////////////////////
	if err := component.InitElasticSearch(config_toml); err != nil {
		return err
	}
	// //////////////////////
	if err := component.InitEcsUploader(config_toml); err != nil {
		return err
	}
	// ///////////////////////
	if err := component.InitLevelDB(config_toml); err != nil {
		return err
	}
	// ///////////////////////
	if err := component.InitSmtpMail(config_toml); err != nil {
		return err
	}
	// ///////////////////////
	if err := component.InitSpr(config_toml); err != nil {
		return err
	}

	// ///////////////////////
	if err := component.InitGeneralCounter(config_toml); err != nil {
		return err
	}

	// ///////////////////////
	// if err := component.InitSqlite(config_toml); err != nil {
	// return err
	// }

	// ////////////////////
	if err := component.InitGeoIp(config_toml); err != nil {
		return err
	}

	return nil
}

func InitComponent() {
	err := InitComponent_()
	if err != nil {
		basic.Logger.Fatalln(err)
	}
}
