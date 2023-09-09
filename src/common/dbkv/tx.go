package dbkv

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/meson-network/bsc-data-file-utils/plugin/redis_plugin"
	"github.com/meson-network/bsc-data-file-utils/src/common/smart_cache"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const DBKV_CACHE_TIME_SECS = 300

func SetDBKV_Str(tx *gorm.DB, keystr string, value string, description string) error {
	tx_result := tx.Table(TABLE_NAME_DBKV).Clauses(clause.OnConflict{UpdateAll: true}).Create(&DBKVModel{Key: keystr, Value: value, Description: description})
	if tx_result.Error != nil {
		return tx_result.Error
	}
	if tx_result.RowsAffected == 0 {
		return errors.New("0 row affected")
	}

	QueryDBKV(tx, nil, &[]string{keystr}, false, true)
	return nil
}

func SetDBKV_UInt64(tx *gorm.DB, keystr string, value uint64, description string) error {
	return SetDBKV_Str(tx, keystr, strconv.FormatUint(value, 10), description)
}

func SetDBKV_Int64(tx *gorm.DB, keystr string, value int64, description string) error {
	return SetDBKV_Str(tx, keystr, strconv.FormatInt(value, 10), description)
}

func SetDBKV_Int32(tx *gorm.DB, keystr string, value int32, description string) error {
	return SetDBKV_Str(tx, keystr, strconv.FormatInt(int64(value), 10), description)
}

func SetDBKV_Int(tx *gorm.DB, keystr string, value int, description string) error {
	return SetDBKV_Str(tx, keystr, strconv.Itoa(value), description)
}

func SetDBKV_Bool(tx *gorm.DB, keystr string, value bool, description string) error {
	if value {
		return SetDBKV_Str(tx, keystr, "true", description)
	} else {
		return SetDBKV_Str(tx, keystr, "false", description)
	}
}

func SetDBKV_Float32(tx *gorm.DB, keystr string, value float32, description string) error {
	return SetDBKV_Str(tx, keystr, fmt.Sprintf("%f", value), description)
}

func SetDBKV_Float64(tx *gorm.DB, keystr string, value float64, description string) error {
	return SetDBKV_Str(tx, keystr, fmt.Sprintf("%f", value), description)
}

func DeleteDBKV_Key(tx *gorm.DB, keystr string) error {
	if err := tx.Table(TABLE_NAME_DBKV).Where(" `key` = ?", keystr).Delete(&DBKVModel{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteDBKV_Id(tx *gorm.DB, id int64) error {
	if err := tx.Table(TABLE_NAME_DBKV).Where(" `id` = ?", id).Delete(&DBKVModel{}).Error; err != nil {
		return err
	}
	return nil
}

func GetDBKV_Bool(tx *gorm.DB, key string) (bool, error) {
	result, err := GetDBKV(tx, nil, &key, true, true)
	if err != nil {
		return false, err
	}

	bool_result, bool_err := result.ToBool()
	if bool_err != nil {
		return false, bool_err
	}
	return bool_result, nil
}

func GetDBKV_Str(tx *gorm.DB, key string) (string, error) {
	result, err := GetDBKV(tx, nil, &key, true, true)
	if err != nil {
		return "", err
	}
	return result.ToString(), nil
}

func GetDBKV_Int(tx *gorm.DB, key string) (int, error) {
	result, err := GetDBKV(tx, nil, &key, true, true)
	if err != nil {
		return 0, err
	}

	int_result, int_err := result.ToInt()
	if int_err != nil {
		return 0, int_err
	}
	return int_result, nil
}

func GetDBKV_Int32(tx *gorm.DB, key string) (int32, error) {
	result, err := GetDBKV(tx, nil, &key, true, true)
	if err != nil {
		return 0, err
	}

	int32_result, int32_err := result.ToInt32()
	if int32_err != nil {
		return 0, int32_err
	}

	return int32_result, nil
}

func GetDBKV_Int64(tx *gorm.DB, key string) (int64, error) {
	result, err := GetDBKV(tx, nil, &key, true, true)
	if err != nil {
		return 0, err
	}

	int64_result, int64_err := result.ToInt64()
	if int64_err != nil {
		return 0, int64_err
	}

	return int64_result, nil
}

func GetDBKV_UInt64(tx *gorm.DB, key string) (uint64, error) {
	result, err := GetDBKV(tx, nil, &key, true, true)
	if err != nil {
		return 0, err
	}

	uint64_result, uint64_err := result.ToUInt64()
	if uint64_err != nil {
		return 0, uint64_err
	}

	return uint64_result, nil
}

func GetDBKV_Float32(tx *gorm.DB, key string) (float32, error) {
	result, err := GetDBKV(tx, nil, &key, true, true)
	if err != nil {
		return 0, err
	}

	float32_result, float32_err := result.ToFloat32()
	if float32_err != nil {
		return 0, float32_err
	}

	return float32_result, nil
}

func GetDBKV_Float64(tx *gorm.DB, key string) (float64, error) {
	result, err := GetDBKV(tx, nil, &key, true, true)
	if err != nil {
		return 0, err
	}

	float64_result, float64_err := result.ToFloat64()
	if float64_err != nil {
		return 0, float64_err
	}

	return float64_result, nil
}

// when error !=nil ,result is nil
// when error == nil , result may be nil or not nil
func GetDBKV(tx *gorm.DB, id *int64, key *string, fromCache bool, updateCache bool) (*DBKVModel, error) {

	var result *DBKVQueryResults
	var err error

	if key == nil {
		result, err = QueryDBKV(tx, id, nil, fromCache, updateCache)
	} else {
		result, err = QueryDBKV(tx, id, &[]string{*key}, fromCache, updateCache)
	}

	if err != nil {
		return nil, err
	}

	if result == nil || result.TotalCount == 0 || len(result.Kv) == 0 {
		return nil, errors.New("nil error")
	}

	return result.Kv[0], nil
}

type DBKVQueryResults struct {
	Kv         []*DBKVModel
	TotalCount int64
}

var query_dbkv_cache_ttl = &smart_cache.QueryCacheTTL{
	Redis_ttl_secs: DBKV_CACHE_TIME_SECS,
	Ref_ttl_secs:   30,
}

func QueryDBKV(tx *gorm.DB, id *int64, keys *[]string, fromCache bool, updateCache bool) (*DBKVQueryResults, error) {

	//gen_key
	ck := smart_cache.NewConnectKey("dbkv")
	ck.C_Int64_Ptr("id", id).C_Str_Array_Ptr("keys", keys)

	key := redis_plugin.GetInstance().GenKey(ck.String())

	/////
	resultHolderAlloc := func() *smart_cache.QueryResult {

		return &smart_cache.QueryResult{
			Result_holder: &DBKVQueryResults{
				Kv:         []*DBKVModel{},
				TotalCount: 0,
			},
			Found:   false,
			Has_err: false,
			Err_str: "",
		}
	}

	/////
	query := func(resultHolder *smart_cache.QueryResult) *smart_cache.QueryCacheTTL {

		queryResults := resultHolder.Result_holder.(*DBKVQueryResults)

		query := tx.Table(TABLE_NAME_DBKV)
		if id != nil {
			query.Where("id = ?", *id)
		}
		if keys != nil {
			query.Where(TABLE_NAME_DBKV+".key IN ?", *keys)
		}

		query.Count(&queryResults.TotalCount)

		err := query.Find(&queryResults.Kv).Error
		if err != nil {
			resultHolder.Has_err = true
			resultHolder.Err_str = err.Error()
			return smart_cache.SlowQueryTTL_ERR
		}

		if queryResults.TotalCount == 0 {
			return smart_cache.SlowQueryTTL_NOT_FOUND
		}

		resultHolder.Found = true
		return query_dbkv_cache_ttl

	}

	s_query := &smart_cache.SlowQuery{
		CacheTTL: query_dbkv_cache_ttl,
		Query:    query,
	}

	/////
	sq_result := smart_cache.SmartQueryCacheSlow(key, fromCache, updateCache, "DBKV Query", resultHolderAlloc, s_query)

	/////
	if sq_result.Has_err {
		return nil, errors.New(sq_result.Err_str)
	} else {
		return sq_result.Result_holder.(*DBKVQueryResults), nil
	}

}
