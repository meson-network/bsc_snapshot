package component

import (
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/basic/config"
	"github.com/meson-network/bsc-data-file-utils/plugin/redis_plugin"
)

func InitRedis(toml_conf *config.TomlConfig) error {

	if toml_conf.Redis.Enable {
		redis_conf := redis_plugin.Config{
			Address:   toml_conf.Redis.Host,
			UserName:  toml_conf.Redis.Username,
			Password:  toml_conf.Redis.Password,
			Port:      toml_conf.Redis.Port,
			KeyPrefix: toml_conf.Redis.Prefix,
			UseTLS:    toml_conf.Redis.Use_tls,
		}

		basic.Logger.Infoln("Init redis plugin with config:", redis_conf)
		if err := redis_plugin.Init(&redis_conf); err == nil {
			basic.Logger.Infoln("### InitRedis success")
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}
