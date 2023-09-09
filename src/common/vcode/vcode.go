package vcode

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/coreservice-io/utils/rand_util"
	goredis "github.com/go-redis/redis/v8"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/plugin/redis_plugin"
)

const redis_vcode_prefix = "vcode"

// send vcode to user
func GenVCode(vCodeKey string, vCodeLen int, vCodeExpireSecs int) (string, error) {
	if vCodeLen <= 0 {
		return "", nil
	}
	key := redis_plugin.GetInstance().GenKey(redis_vcode_prefix, vCodeKey)
	code, _ := redis_plugin.GetInstance().Get(context.Background(), key).Result()
	if len(code) != vCodeLen {
		code = rand_util.GenRandStr(vCodeLen)
	}
	_, err := redis_plugin.GetInstance().Set(context.Background(), key, code, time.Duration(vCodeExpireSecs)*time.Second).Result()
	if err != nil {
		basic.Logger.Errorln("GenVCode set email vcode to redis error", "err", err)
		return "", errors.New("set email vcode error")
	}

	basic.Logger.Debugln("vcode", "code", code, "vCodeKey", vCodeKey)
	return code, nil
}

func ValidateVCode(vCodeKey string, code string) bool {
	//incase user have Cap or whitespace
	code = strings.ToLower(code)
	code = strings.TrimSpace(code)

	key := redis_plugin.GetInstance().GenKey(redis_vcode_prefix, vCodeKey)
	value, err := redis_plugin.GetInstance().Get(context.Background(), key).Result()
	if err == goredis.Nil {
		return false
	} else if err != nil {
		basic.Logger.Debugln("ValidateVCode from redis err", "err", err, "vCodeKey", vCodeKey)
		return false
	}
	redis_plugin.GetInstance().Del(context.Background(), key)
	return value == code
}

func ValidateVCodeWithLen(vCodeKey string, code string, required_len int) bool {
	code = strings.ToLower(code)
	code = strings.TrimSpace(code)
	if len(code) != required_len {
		return false
	}
	return ValidateVCode(vCodeKey, code)
}
