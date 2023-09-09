package component

import (
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/basic/config"
	"github.com/meson-network/bsc-data-file-utils/plugin/ecs_plugin"
)

func InitElasticSearch(toml_conf *config.TomlConfig) error {

	if toml_conf.Elastic_search.Enable {

		ecs_conf := ecs_plugin.Config{
			Address:  toml_conf.Elastic_search.Host,
			UserName: toml_conf.Elastic_search.Username,
			Password: toml_conf.Elastic_search.Password}

		basic.Logger.Infoln("Init elastic search plugin with config:", ecs_conf)
		if err := ecs_plugin.Init(&ecs_conf); err == nil {
			basic.Logger.Infoln("### InitElasticSearch success")
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}
