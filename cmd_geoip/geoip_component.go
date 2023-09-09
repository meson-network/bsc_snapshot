package cmd_geoip

import (
	"errors"

	"github.com/coreservice-io/geo_ip/lib"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/basic/config"
)

func StartGeoIpComponent() (lib.GeoIpInterface, error) {
	toml_conf := config.Get_config().Toml_config
	/////////////////////////
	if !toml_conf.Geo_ip.Enable {
		basic.Logger.Fatalln("geoip not enabled in config")
	}

	if toml_conf.Geo_ip.Update_key == "" {
		return nil, errors.New("geoip update_key error:" + toml_conf.Geo_ip.Update_key)
	}

	if toml_conf.Geo_ip.Dataset_version == "" {
		return nil, errors.New("geoip dataset_version error:" + toml_conf.Geo_ip.Dataset_version)
	}

	_, ds_f_abs_exist, _ := basic.PathExist(toml_conf.Geo_ip.Dataset_folder)
	if !ds_f_abs_exist {
		return nil, errors.New("geo_ip Dataset_folder path error," + toml_conf.Geo_ip.Dataset_folder)
	}

	client, err := lib.NewClient(toml_conf.Geo_ip.Update_key, toml_conf.Geo_ip.Dataset_version,
		toml_conf.Geo_ip.Dataset_folder, true, func(s string) {
			basic.Logger.Debugln(s)
		}, func(s string) {
			basic.Logger.Errorln(s)
		})

	if err != nil {
		return nil, err
	}

	return client, nil

}
