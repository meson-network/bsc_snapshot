package smart_cache

import (
	"fmt"
	"strconv"

	"github.com/coreservice-io/utils/hash_util"
)

type ConnectKey struct {
	Key string
}

func NewConnectKey(baseKey string) *ConnectKey {
	return &ConnectKey{
		Key: baseKey,
	}
}

func (k *ConnectKey) String() string {
	return k.Key
}

func (k *ConnectKey) C_Str(key string) *ConnectKey {
	k.Key = k.Key + ":" + key
	return k
}

func (k *ConnectKey) C_Str_Ptr(keyname string, key_ptr *string) *ConnectKey {
	if key_ptr != nil {
		k.C_Str(*key_ptr)
	} else {
		k.Key = k.Key + ":" + keyname
	}
	return k
}

func (k *ConnectKey) C_Int(key int) *ConnectKey {
	k.Key = k.Key + ":" + strconv.Itoa(key)
	return k
}

func (k *ConnectKey) C_Int_Ptr(keyname string, key_ptr *int) *ConnectKey {
	if key_ptr != nil {
		k.C_Int(*key_ptr)
	} else {
		k.Key = k.Key + ":" + keyname
	}
	return k
}

func (k *ConnectKey) C_Int64(key int64) *ConnectKey {
	k.Key = k.Key + ":" + strconv.FormatInt(key, 10)
	return k
}

func (k *ConnectKey) C_Int64_Ptr(keyname string, key_ptr *int64) *ConnectKey {
	if key_ptr != nil {
		k.C_Int64(*key_ptr)
	} else {
		k.Key = k.Key + ":" + keyname
	}
	return k
}

func (k *ConnectKey) C_UInt(key uint) *ConnectKey {
	k.Key = k.Key + ":" + strconv.FormatUint(uint64(key), 10)
	return k
}

func (k *ConnectKey) C_UInt_Ptr(keyname string, key_ptr *uint) *ConnectKey {
	if key_ptr != nil {
		k.C_UInt(*key_ptr)
	} else {
		k.Key = k.Key + ":" + keyname
	}
	return k
}

func (k *ConnectKey) C_UInt64(key uint64) *ConnectKey {
	k.Key = k.Key + ":" + strconv.FormatUint(key, 10)
	return k
}

func (k *ConnectKey) C_UInt64_Ptr(keyname string, key_ptr *uint64) *ConnectKey {
	if key_ptr != nil {
		k.C_UInt64(*key_ptr)
	} else {
		k.Key = k.Key + ":" + keyname
	}
	return k
}

func (k *ConnectKey) C_Bool(key bool) *ConnectKey {
	if key {
		k.Key = k.Key + ":" + "true"
	} else {
		k.Key = k.Key + ":" + "false"
	}
	return k
}

func (k *ConnectKey) C_Bool_Ptr(keyname string, key_ptr *bool) *ConnectKey {
	if key_ptr != nil {
		k.C_Bool(*key_ptr)
	} else {
		k.Key = k.Key + ":" + keyname
	}
	return k
}

func (k *ConnectKey) C_Float32(key float32) *ConnectKey {
	k.Key = k.Key + ":" + fmt.Sprintf("%f", key)
	return k
}

func (k *ConnectKey) C_Float32_Ptr(keyname string, key_ptr *float32) *ConnectKey {
	if key_ptr != nil {
		k.C_Float32(*key_ptr)
	} else {
		k.Key = k.Key + ":" + keyname
	}
	return k
}

func (k *ConnectKey) C_Float64(key float64) *ConnectKey {
	k.Key = k.Key + ":" + fmt.Sprintf("%f", key)
	return k
}

func (k *ConnectKey) C_Float64_Ptr(keyname string, key_ptr *float64) *ConnectKey {
	if key_ptr != nil {
		k.C_Float64(*key_ptr)
	} else {
		k.Key = k.Key + ":" + keyname
	}
	return k
}

func (k *ConnectKey) C_Int_Array(key []int) *ConnectKey {

	key_str := "array_" + strconv.Itoa(len(key))
	for _, v := range key {
		key_str = key_str + "_" + strconv.Itoa(v)
	}

	if len(key_str) > 10 {
		key_str = hash_util.CRC32HashString(key_str)
	}

	k.Key = k.Key + ":" + key_str
	return k
}

func (k *ConnectKey) C_Int_Array_Ptr(keyname string, key_ptr *[]int) *ConnectKey {
	if key_ptr != nil {
		k.C_Int_Array(*key_ptr)
	} else {
		k.Key = k.Key + ":" + keyname
	}
	return k
}

func (k *ConnectKey) C_Int64_Array(key []int64) *ConnectKey {

	key_str := "array_" + strconv.Itoa(len(key))
	for _, v := range key {
		key_str = key_str + "_" + strconv.FormatInt(v, 10)
	}

	if len(key_str) > 10 {
		key_str = hash_util.CRC32HashString(key_str)
	}

	k.Key = k.Key + ":" + key_str
	return k
}

func (k *ConnectKey) C_Int64_Array_Ptr(keyname string, key_ptr *[]int64) *ConnectKey {
	if key_ptr != nil {
		k.C_Int64_Array(*key_ptr)
	} else {
		k.Key = k.Key + ":" + keyname
	}
	return k
}

func (k *ConnectKey) C_Str_Array(key []string) *ConnectKey {

	key_str := "array_" + strconv.Itoa(len(key))
	for _, v := range key {
		key_str = key_str + "_" + v
	}

	if len(key_str) > 10 {
		key_str = hash_util.CRC32HashString(key_str)
	}

	k.Key = k.Key + ":" + key_str
	return k
}

func (k *ConnectKey) C_Str_Array_Ptr(keyname string, key_ptr *[]string) *ConnectKey {
	if key_ptr != nil {
		k.C_Str_Array(*key_ptr)
	} else {
		k.Key = k.Key + ":" + keyname
	}
	return k
}
