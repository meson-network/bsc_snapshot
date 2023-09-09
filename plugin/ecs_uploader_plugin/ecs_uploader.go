package ecs_uploader_plugin

import (
	"fmt"

	"github.com/coreservice-io/ecs_uploader/uploader"
	"github.com/coreservice-io/log"
	"github.com/meson-network/bsc-data-file-utils/basic"
)

var instanceMap = map[string]*uploader.Uploader{}

func GetInstance() *uploader.Uploader {
	return GetInstance_("default")
}

func GetInstance_(name string) *uploader.Uploader {
	ecs_uploader_i := instanceMap[name]
	if ecs_uploader_i == nil {
		basic.Logger.Errorln(name + " ecs_uploader plugin null")
	}
	return ecs_uploader_i
}

/*
elasticSearchAddr
elasticSearchUserName
elasticSearchPassword
*/
type Config struct {
	Address  string
	UserName string
	Password string
}

func Init(esConfig *Config, logger log.Logger) error {
	return Init_("default", esConfig, logger)
}

// Init a new instance.
// If only need one instance, use empty name "". Use GetDefaultInstance() to get.
// If you need several instance, run Init() with different <name>. Use GetInstance(<name>) to get.
func Init_(name string, esConfig *Config, logger log.Logger) error {
	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("ecs_uploader instance <%s> has already been initialized", name)
	}

	es, err := uploader.New(esConfig.Address, esConfig.UserName, esConfig.Password)
	if err != nil {
		return err
	}

	es.SetLogger(logger)
	instanceMap[name] = es
	return nil
}
