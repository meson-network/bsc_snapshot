package component

import (
	"errors"

	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/basic/config"
	"github.com/meson-network/bsc-data-file-utils/plugin/geo_ip_plugin"
)

func InitGeoIp(toml_conf *config.TomlConfig) error {

	if toml_conf.Geo_ip.Enable {

		if toml_conf.Geo_ip.Update_key == "" {
			return errors.New("geoip update_key error:" + toml_conf.Geo_ip.Update_key)
		}

		if toml_conf.Geo_ip.Dataset_version == "" {
			return errors.New("geoip dataset_version error:" + toml_conf.Geo_ip.Dataset_version)
		}

		ds_f_abs, ds_f_abs_exist, _ := basic.PathExist(toml_conf.Geo_ip.Dataset_folder)
		if !ds_f_abs_exist {
			return errors.New("geo_ip Dataset_folder path error," + toml_conf.Geo_ip.Dataset_folder)
		}

		basic.Logger.Debugln("Init geo_ip plugin with dataset_version:"+toml_conf.Geo_ip.Dataset_version+" and Dataset_folder:", ds_f_abs)

		if err := geo_ip_plugin.Init(toml_conf.Geo_ip.Update_key, toml_conf.Geo_ip.Dataset_version, ds_f_abs, func(s string) {
			basic.Logger.Debugln(s)
		}, func(s string) {
			basic.Logger.Errorln(s)
		}); err == nil {
			basic.Logger.Debugln("### Init geo_ip success")
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}
