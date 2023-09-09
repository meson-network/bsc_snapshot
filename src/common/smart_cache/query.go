package smart_cache

import (
	"context"
	"errors"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/plugin/redis_plugin"
	"github.com/meson-network/bsc-data-file-utils/plugin/reference_plugin"
)

type QueryCacheTTL struct {
	Redis_ttl_secs int64
	Ref_ttl_secs   int64
}

type QueryResult struct {
	Result_holder interface{}
	Found         bool
	Has_err       bool
	Err_str       string
}

type SlowQuery struct {
	//return default of : redis_ttl_secs, ref_ttl_secs
	CacheTTL *QueryCacheTTL
	//return redis_ttl_secs, ref_ttl_secs,error
	Query func(qResult *QueryResult) *QueryCacheTTL
}

var SlowQueryTTL_Default = &QueryCacheTTL{
	Redis_ttl_secs: 300,
	Ref_ttl_secs:   5,
}

var SlowQueryTTL_NOT_FOUND = &QueryCacheTTL{
	Redis_ttl_secs: 30,
	Ref_ttl_secs:   5,
}

var SlowQueryTTL_ERR = &QueryCacheTTL{
	Redis_ttl_secs: 5,
	Ref_ttl_secs:   5,
}

/////////////////

type smartCacheRefElement struct {
	Obj        interface{}
	Token_chan chan struct{}
}

var lockMap sync.Map

func CheckLimitOffset(limit *int, offset *int) error {

	if offset != nil && limit == nil {
		return errors.New("err: limit missing")
	}

	if limit != nil && *limit <= 0 {
		return errors.New("err: limit >0 required")
	}

	if offset != nil && *offset < 0 {
		return errors.New("err: offset >=0 required")
	}

	return nil
}

func SmartQueryCacheSlow(key string, fromCache bool, updateCache bool, queryDescription string,
	resultHolderAlloc func() *QueryResult, slowQuery *SlowQuery) *QueryResult {

	if fromCache {
		// try to get from reference
		refElement, to_update_ref := refGet(reference_plugin.GetInstance(), key)

		if refElement != nil && !to_update_ref { // 1. ref exist and no need to update
			basic.Logger.Debugln(queryDescription + " SmartQueryCacheSlow hit from reference")
			return refElement.Obj.(*QueryResult)

		} else if refElement != nil && to_update_ref { //2. ref exist and need update
			select {
			case <-refElement.Token_chan: //get run token
				go func() {
					defer func() {
						refElement.Token_chan <- struct{}{} // release run token
					}()

					resultHolder := resultHolderAlloc()
					// get from redis
					basic.Logger.Debugln(queryDescription, " SmartQueryCacheSlow try from redis")
					redis_get_err := redisGet(context.Background(), redis_plugin.GetInstance().ClusterClient, key, resultHolder)

					if redis_get_err == nil {
						// exist in redis
						// ref update
						refElement.Obj = resultHolder
						refSetTTL(reference_plugin.GetInstance(), key, refElement, slowQuery.CacheTTL.Ref_ttl_secs+REF_TTL_DELAY_SECS)
					} else {
						if redis_get_err == redis.Nil {
							//redis:no error just not found
							//try from origin (example form db)
							query_ttl := slowQuery.Query(resultHolder)
							if !resultHolder.Has_err {
								refElement.Obj = resultHolder
								refSetTTL(reference_plugin.GetInstance(), key, refElement, query_ttl.Ref_ttl_secs+REF_TTL_DELAY_SECS)
								redisSet(context.Background(), redis_plugin.GetInstance().ClusterClient, key, resultHolder, query_ttl.Redis_ttl_secs)
							} else {
								//db error ,redis no record , just keep old local ttl longer
								//nothing to do but keep old local reference longer with some secs ,e.g: retry redis after some secs
								refSetTTL(reference_plugin.GetInstance(), key, refElement, query_ttl.Ref_ttl_secs+REF_TTL_DELAY_SECS)
							}
						} else {
							//redis:other err
							//nothing to do but keep old local reference longer with some secs ,e.g: retry redis after some secs
							refSetTTL(reference_plugin.GetInstance(), key, refElement, QUERY_ERR_REF_TTL_SECS+REF_TTL_DELAY_SECS)
						}
					}
				}()
				return refElement.Obj.(*QueryResult)
			default:
				return refElement.Obj.(*QueryResult)
			}
		} else {
			//3. ref not exist
			lc, loaded := lockMap.LoadOrStore(key, make(chan struct{}, 1))
			if loaded { //most query enter this filed
				<-lc.(chan struct{}) //all processes unblock when chan closed
				refElement, _ := refGet(reference_plugin.GetInstance(), key)
				if refElement != nil {
					return refElement.Obj.(*QueryResult)
				} else {
					basic.Logger.Errorln("SmartQueryCacheSlow ref nothing found which may caused by too short ref ttl which seems all most impossible")
					return &QueryResult{
						Result_holder: resultHolderAlloc(),
						Found:         false,
						Has_err:       true,
						Err_str:       "SmartQueryCacheSlow ref nothing found,all most impossible",
					}
				}

			} else {
				// only 1 query enter below
				// get from redis
				basic.Logger.Debugln(queryDescription, " SmartQueryCacheSlow try from redis")
				//////
				resultHolder := resultHolderAlloc()
				tokenChan := make(chan struct{}, 1)
				tokenChan <- struct{}{}
				ele := &smartCacheRefElement{
					Obj:        resultHolder,
					Token_chan: tokenChan, // a new chan
				}
				///////
				redis_get_err := redisGet(context.Background(), redis_plugin.GetInstance().ClusterClient, key, resultHolder)

				if redis_get_err == nil {
					//redis found
					refSetTTL(reference_plugin.GetInstance(), key, ele, slowQuery.CacheTTL.Ref_ttl_secs+REF_TTL_DELAY_SECS)
				} else {
					if redis_get_err == redis.Nil {
						//redis no err but not found
						query_ttl := slowQuery.Query(resultHolder)
						if !resultHolder.Has_err {
							refSetTTL(reference_plugin.GetInstance(), key, ele, query_ttl.Ref_ttl_secs+REF_TTL_DELAY_SECS)
							redisSet(context.Background(), redis_plugin.GetInstance().ClusterClient, key, resultHolder, query_ttl.Redis_ttl_secs)
						} else {
							refSetTTL(reference_plugin.GetInstance(), key, ele, query_ttl.Ref_ttl_secs+REF_TTL_DELAY_SECS)
						}

					} else {
						// redis other error
						// as this is the case that : there is no reference and redis error happens
						// just put an error in the reference
						resultHolder.Has_err = true
						resultHolder.Err_str = redis_get_err.Error()
						refSetTTL(reference_plugin.GetInstance(), key, ele, QUERY_ERR_REF_TTL_SECS+REF_TTL_DELAY_SECS)
					}
				}

				close(lc.(chan struct{}))
				lockMap.Delete(key)
				return resultHolder
			}
		}
	} else {
		// after cache miss ,try from remote database
		resultHolder := resultHolderAlloc()

		s_cache_ttl := slowQuery.Query(resultHolder)

		if updateCache {
			if !resultHolder.Has_err {

				ele, _ := refGet(reference_plugin.GetInstance(), key)
				if ele == nil {
					tokenChan := make(chan struct{}, 1)
					tokenChan <- struct{}{}
					ele = &smartCacheRefElement{
						Obj:        resultHolder,
						Token_chan: tokenChan, // a new chan
					}
				} else {
					ele.Obj = resultHolder
				}

				refSetTTL(reference_plugin.GetInstance(), key, ele, s_cache_ttl.Ref_ttl_secs+REF_TTL_DELAY_SECS)
				redisSet(context.Background(), redis_plugin.GetInstance().ClusterClient, key, resultHolder, s_cache_ttl.Redis_ttl_secs)
			}
			//// else nothing ,as cache should kept old values
		}

		return resultHolder
	}
}

///////////////////////////////////////////////////////////////////////////////////

// fastQuery usually is a redis query
func SmartQueryCacheFast(
	key string,
	resultHolderAlloc func() *QueryResult,
	fromRefCache bool,
	updateRefCache bool,
	fastQuery func(resultHolder *QueryResult) int64, //return ref_ttl_secs ,err
	queryDescription string) *QueryResult {

	if fromRefCache {
		// try to get from reference
		refElement, to_update_ref := refGet(reference_plugin.GetInstance(), key)

		if refElement != nil && !to_update_ref { // 1. ref exist and no need to update

			basic.Logger.Debugln(queryDescription + " SmartQueryCacheFast hit from reference")
			return refElement.Obj.(*QueryResult)

		} else if refElement != nil && to_update_ref { //2. ref exist and need update
			select {
			case <-refElement.Token_chan: //get run token
				go func() {
					defer func() {
						refElement.Token_chan <- struct{}{} // release run token
					}()

					//try from origin
					resultHolder := resultHolderAlloc()
					refCacheTTLSecs := fastQuery(resultHolder)

					if !resultHolder.Has_err {
						refElement.Obj = resultHolder
						refSetTTL(reference_plugin.GetInstance(), key, refElement, refCacheTTLSecs+REF_TTL_DELAY_SECS)
					} else {
						//nothing to do but keep old local ref some secs longer
						refSetTTL(reference_plugin.GetInstance(), key, refElement, QUERY_ERR_REF_TTL_SECS+REF_TTL_DELAY_SECS)
					}
				}()
				return refElement.Obj.(*QueryResult)
			default:
				return refElement.Obj.(*QueryResult)
			}
		} else {

			//3. ref not exist
			lc, loaded := lockMap.LoadOrStore(key, make(chan struct{}, 1))
			if loaded { //most query enter this filed
				<-lc.(chan struct{}) //unblock when chan closed

				refElement, _ := refGet(reference_plugin.GetInstance(), key)
				if refElement != nil {
					return refElement.Obj.(*QueryResult)
				} else {
					basic.Logger.Errorln("SmartQueryCacheFast ref nothing found which may caused by too short ref ttl which seems all most impossible")
					return &QueryResult{
						Result_holder: resultHolderAlloc(),
						Found:         false,
						Has_err:       true,
						Err_str:       "SmartQueryCacheFast ref nothing found,all most impossible",
					}
				}

			} else {
				// only 1 query enter below
				//try from origin
				resultHolder := resultHolderAlloc()
				tokenChan := make(chan struct{}, 1)
				tokenChan <- struct{}{}
				ele := &smartCacheRefElement{
					Obj:        resultHolder,
					Token_chan: tokenChan, // a new chan
				}

				refCacheTTLSecs := fastQuery(resultHolder)
				refSetTTL(reference_plugin.GetInstance(), key, ele, refCacheTTLSecs+REF_TTL_DELAY_SECS)

				close(lc.(chan struct{})) //just close to unlock chan
				lockMap.Delete(key)
				return resultHolder
			}
		}
	} else {
		//after cache miss ,try from remote database
		basic.Logger.Debugln(queryDescription, " SmartQueryCacheFast try from fast query")
		resultHolder := resultHolderAlloc()

		f_cache_ttl_secs := fastQuery(resultHolder)

		if updateRefCache {

			if !resultHolder.Has_err {

				ele, _ := refGet(reference_plugin.GetInstance(), key)
				if ele == nil {
					tokenChan := make(chan struct{}, 1)
					tokenChan <- struct{}{}
					ele = &smartCacheRefElement{
						Obj:        resultHolder,
						Token_chan: tokenChan, // a new chan
					}
				} else {
					ele.Obj = resultHolder
				}

				refSetTTL(reference_plugin.GetInstance(), key, ele, f_cache_ttl_secs+REF_TTL_DELAY_SECS)
			}
			//nothing changes
		}
		return resultHolder
	}
}
