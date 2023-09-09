package redis_pipeline

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	operation_Set              = "Set"
	operation_ZAdd             = "ZAdd"
	operation_ZAddNX           = "ZAddNX"
	operation_HSet             = "HSet"
	operation_Expire           = "Expire"
	operation_ZRemRangeByScore = "ZRemRangeByScore"
)

//don't use pipeline for high-safety senario
//as redis pipeline may have data lost problem
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	redisCmd := &PipelineCmd{
		Ctx:       ctx,
		Operation: operation_Set,
		Key:       key,
		Args:      []interface{}{value, expiration},
	}
	cmdListChannel <- redisCmd
}

//don't use pipeline for high-safety senario
//as redis pipeline may have data lost problem
func ZAdd(ctx context.Context, key string, members ...*redis.Z) {
	redisCmd := &PipelineCmd{
		Ctx:       ctx,
		Operation: operation_ZAdd,
		Key:       key,
		Args:      []interface{}{},
	}
	for _, v := range members {
		redisCmd.Args = append(redisCmd.Args, v)
	}
	cmdListChannel <- redisCmd
}

//don't use pipeline for high-safety senario
//as redis pipeline may have data lost problem
func ZAddNX(ctx context.Context, key string, members ...*redis.Z) {
	redisCmd := &PipelineCmd{
		Ctx:       ctx,
		Operation: operation_ZAddNX,
		Key:       key,
		Args:      []interface{}{},
	}
	for _, v := range members {
		redisCmd.Args = append(redisCmd.Args, v)
	}
	cmdListChannel <- redisCmd
}

//don't use pipeline for high-safety senario
//as redis pipeline may have data lost problem
func HSet(ctx context.Context, key string, values ...interface{}) {
	redisCmd := &PipelineCmd{
		Ctx:       ctx,
		Operation: operation_HSet,
		Key:       key,
		Args:      []interface{}{},
	}
	redisCmd.Args = append(redisCmd.Args, values...)
	cmdListChannel <- redisCmd
}

//don't use pipeline for high-safety senario
//as redis pipeline may have data lost problem
func Expire(ctx context.Context, key string, expiration time.Duration) {
	redisCmd := &PipelineCmd{
		Ctx:       ctx,
		Operation: operation_Expire,
		Key:       key,
		Args:      []interface{}{expiration},
	}
	cmdListChannel <- redisCmd
}

//don't use pipeline for high-safety senario
//as redis pipeline may have data lost problem
func ZRemRangeByScore(ctx context.Context, key, min, max string) {
	redisCmd := &PipelineCmd{
		Ctx:       ctx,
		Operation: operation_ZRemRangeByScore,
		Key:       key,
		Args:      []interface{}{min, max},
	}
	cmdListChannel <- redisCmd
}
