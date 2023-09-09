package spr_plugin

import (
	"fmt"

	"github.com/coreservice-io/log"
	"github.com/coreservice-io/redis_spr"
	"github.com/meson-network/bsc-data-file-utils/basic"
)

var instanceMap = map[string]*redis_spr.SprJobMgr{}

func GetInstance() *redis_spr.SprJobMgr {
	return GetInstance_("default")
}

func GetInstance_(name string) *redis_spr.SprJobMgr {
	spr_i := instanceMap[name]
	if spr_i == nil {
		basic.Logger.Errorln(name + " spr plugin null")
	}
	return spr_i
}

func Init(redisConfig *redis_spr.RedisConfig, logger log.Logger) error {
	return Init_("default", redisConfig, logger)
}

// Init a new instance.
//
//	If only need one instance, use empty name "". Use GetDefaultInstance() to get.
//	If you need several instance, run Init() with different <name>. Use GetInstance(<name>) to get.
func Init_(name string, redisConfig *redis_spr.RedisConfig, logger log.Logger) error {
	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("spr instance <%s> has already been initialized", name)
	}

	if redisConfig.Addr == "" {
		redisConfig.Addr = "127.0.0.1"
	}
	if redisConfig.Port == 0 {
		redisConfig.Port = 6379
	}
	//////// ini spr job //////////////////////

	spr, err := redis_spr.New(redis_spr.RedisConfig{
		Addr:     redisConfig.Addr,
		Port:     redisConfig.Port,
		Password: redisConfig.Password,
		UserName: redisConfig.UserName,
		Prefix:   redisConfig.Prefix,
		UseTLS:   redisConfig.UseTLS,
	})

	if err != nil {
		return err
	}

	spr.SetLogger(logger)

	instanceMap[name] = spr

	return nil
}
