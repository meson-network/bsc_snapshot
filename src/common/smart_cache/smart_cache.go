package smart_cache

import (
	"context"
	"time"

	"github.com/coreservice-io/reference"
	"github.com/go-redis/redis/v8"
	"github.com/meson-network/bsc-data-file-utils/src/common/json"
)

const REF_TTL_DELAY_SECS = 5 //add REF_TTL_DELAY_SECS to local ref when set

const QUERY_ERR_REF_TTL_SECS = 5   //the ref tempory store secs when query error
const QUERY_ERR_REDIS_TTL_SECS = 5 //the redis tempory store secs when query error

const QUERY_NIL_REF_TTL_SECS_DEFAULT = 5    //the ref tempory store secs when query nil
const QUERY_NIL_REDIS_TTL_SECS_DEFAULT = 15 //the redis tempory store secs when query nil

func refGet(localRef *reference.Reference, keystr string) (result *smartCacheRefElement, to_update bool) {
	refElement, ttl := localRef.Get(keystr)
	if refElement == nil {
		return nil, true
	}

	if ttl <= REF_TTL_DELAY_SECS {
		return refElement.(*smartCacheRefElement), true
	} else {
		return refElement.(*smartCacheRefElement), false
	}
}

func refSetTTL(localRef *reference.Reference, keystr string, element *smartCacheRefElement, ref_ttl_second int64) error {
	return localRef.Set(keystr, element, ref_ttl_second)
}

// return to_update bool
func redisGet(ctx context.Context, Redis *redis.ClusterClient, keystr string, q_result *QueryResult) error {

	scmd := Redis.Get(ctx, keystr)
	r_bytes, err := scmd.Bytes()

	if err != nil {
		return err
	}

	return json.Unmarshal(r_bytes, q_result)
}

func redisSet(ctx context.Context, Redis *redis.ClusterClient, keystr string, q_result *QueryResult, redis_ttl_second int64) error {
	v_json, err := json.Marshal(q_result)
	if err != nil {
		return err
	}
	return Redis.Set(ctx, keystr, v_json, time.Duration(redis_ttl_second)*time.Second).Err()
}
