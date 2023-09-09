package redis_plugin

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/meson-network/bsc-data-file-utils/basic"
)

var instanceMap = map[string]*RedisClient{}

type RedisClient struct {
	KeyPrefix string
	*redis.ClusterClient
}

func GetInstance() *RedisClient {
	return GetInstance_("default")
}

func GetInstance_(name string) *RedisClient {
	redis_i := instanceMap[name]
	if redis_i == nil {
		basic.Logger.Errorln(name + " redis plugin null")
	}
	return redis_i
}

type Config struct {
	Address   string
	UserName  string
	Password  string
	Port      int
	KeyPrefix string
	UseTLS    bool
}

func Init(redisConfig *Config) error {
	return Init_("default", redisConfig)
}

// Init a new instance.
//
//	If only need one instance, use empty name "". Use GetDefaultInstance() to get.
//	If you need several instance, run Init() with different <name>. Use GetInstance(<name>) to get.
func Init_(name string, redisConfig *Config) error {
	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("redis instance <%s> has already been initialized", name)
	}

	if redisConfig.Address == "" {
		redisConfig.Address = "127.0.0.1"
	}
	if redisConfig.Port == 0 {
		redisConfig.Port = 6379
	}

	config := &redis.ClusterOptions{
		Addrs:    []string{redisConfig.Address + ":" + strconv.Itoa(redisConfig.Port)},
		Username: redisConfig.UserName,
		Password: redisConfig.Password,
	}
	if redisConfig.UseTLS {
		config.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	r := redis.NewClusterClient(config)

	_, err := r.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	prefix := redisConfig.KeyPrefix
	if prefix != "" && !strings.HasSuffix(prefix, ":") {
		prefix = prefix + ":"
	}

	instanceMap[name] = &RedisClient{
		prefix,
		r,
	}
	return nil
}

func (rc *RedisClient) GenKey(keys ...string) string {
	if len(keys) == 0 {
		return rc.KeyPrefix
	}
	return rc.KeyPrefix + strings.Join(keys, ":")
}
