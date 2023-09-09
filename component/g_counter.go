package component

import (
	general_counter "github.com/coreservice-io/general-counter"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/basic/config"
	"github.com/meson-network/bsc-data-file-utils/plugin/g_counter_plugin"
)

func InitGeneralCounter(toml_conf *config.TomlConfig) error {
	if toml_conf.General_counter.Enable {
		g_conf := &general_counter.GeneralCounterConfig{
			Project_name: toml_conf.General_counter.Project_name,
			Db_config: &general_counter.DBConfig{
				Host:     toml_conf.Db.Host,
				Port:     toml_conf.Db.Port,
				DbName:   toml_conf.Db.Name,
				UserName: toml_conf.Db.Username,
				Password: toml_conf.Db.Password,
			},
			Ecs_config: &general_counter.EcsConfig{
				Address:  toml_conf.Elastic_search.Host,
				UserName: toml_conf.Elastic_search.Username,
				Password: toml_conf.Elastic_search.Password,
			},
			Redis_config: &general_counter.RedisConfig{
				Addr:     toml_conf.Redis.Host,
				UserName: toml_conf.Redis.Username,
				Password: toml_conf.Redis.Password,
				Port:     toml_conf.Redis.Port,
				Prefix:   toml_conf.Redis.Prefix,
				UseTLS:   toml_conf.Redis.Use_tls,
			},
		}

		basic.Logger.Infoln("Init g_counter plugin with config:", g_conf)
		if err := g_counter_plugin.Init(g_conf, basic.Logger); err == nil {
			basic.Logger.Infoln("### Init g_counter success")
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}
