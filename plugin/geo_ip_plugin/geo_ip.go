package geo_ip_plugin

import (
	"fmt"

	"github.com/coreservice-io/geo_ip/lib"
	"github.com/meson-network/bsc-data-file-utils/basic"
)

type GeoIp struct {
	lib.GeoIpInterface
}

var instanceMap = map[string]*GeoIp{}

func GetInstance() *GeoIp {
	return GetInstance_("default")
}

func GetInstance_(name string) *GeoIp {
	geoip := instanceMap[name]
	if geoip == nil {
		basic.Logger.Errorln(name + " geoip plugin null")
	}
	return geoip
}

func Init(update_key string, version string, dataset_folder string, logger func(string), err_logger func(string)) error {
	return Init_("default", update_key, version, dataset_folder, logger, err_logger)
}

func Init_(name string, update_key string, version string, dataset_folder string, logger func(string), err_logger func(string)) error {
	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("ip_geo instance <%s> has already initialized", name)
	}

	ipClient := &GeoIp{}
	// new instance
	client, err := lib.NewClient(update_key, version, dataset_folder, false, logger, err_logger)
	if err != nil {
		return err
	}
	ipClient.GeoIpInterface = client

	instanceMap[name] = ipClient
	return nil
}
