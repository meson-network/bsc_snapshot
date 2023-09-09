package dbkv

import (
	"errors"
	"strconv"
)

const TABLE_NAME_DBKV = "dbkv"

type DBKVModel struct {
	Id               int64  `json:"id" gorm:"type:bigint(20);primaryKey;autoIncrement"`
	Key              string `json:"key" gorm:"type:varchar(512);uniqueIndex"`
	Value            string `json:"value" gorm:"type:longtext;"`
	Description      string `json:"description" gorm:"type:longtext;"`
	Update_unixtime  int64  `json:"update_unixtime" gorm:"autoUpdateTime;type:bigint(20)"`
	Created_unixtime int64  `json:"created_unixtime" gorm:"autoCreateTime;type:bigint(20)"`
}

func (model *DBKVModel) TableName() string {
	return TABLE_NAME_DBKV
}

func (model *DBKVModel) ToString() string {
	return model.Value
}

func (model *DBKVModel) ToBool() (bool, error) {
	switch model.Value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, errors.New("format err, str value:" + model.Value)
	}
}

func (model *DBKVModel) ToInt() (int, error) {
	return strconv.Atoi(model.Value)
}

func (model *DBKVModel) ToInt32() (int32, error) {
	result, err := strconv.ParseInt(model.Value, 10, 32)
	if err != nil {
		return 0, err
	} else {
		return int32(result), nil
	}
}

func (model *DBKVModel) ToInt64() (int64, error) {
	return strconv.ParseInt(model.Value, 10, 64)
}

func (model *DBKVModel) ToUInt64() (uint64, error) {
	return strconv.ParseUint(model.Value, 10, 64)
}

func (model *DBKVModel) ToFloat64() (float64, error) {
	return strconv.ParseFloat(model.Value, 64)
}

func (model *DBKVModel) ToFloat32() (float32, error) {
	result, err := strconv.ParseFloat(model.Value, 32)
	if err != nil {
		return 0, err
	} else {
		return float32(result), nil
	}
}
