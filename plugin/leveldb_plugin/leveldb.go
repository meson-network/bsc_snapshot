package leveldb_plugin

import (
	"fmt"

	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/syndtr/goleveldb/leveldb"
)

type Config struct {
	Db_folder string
}

var instanceMap = map[string]*leveldb.DB{}

func GetInstance() *leveldb.DB {
	return GetInstance_("default")
}

func GetInstance_(name string) *leveldb.DB {
	ldb := instanceMap[name]
	if ldb == nil {
		basic.Logger.Errorln(name + "level db plugin null")
	}
	return ldb
}

func Init(config *Config) error {
	return Init_("default", config)
}

// Init a new instance.
// If only need one instance, use empty name "". Use GetDefaultInstance() to get.
// If you need several instance, run Init() with different <name>. Use GetInstance(<name>) to get.
func Init_(name string, config *Config) error {
	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("*leveldb.DB instance <%s> has already been initialized", name)
	}

	////
	db, err := leveldb.OpenFile(config.Db_folder, nil)
	if err != nil {
		return err
	}

	instanceMap[name] = db
	return nil
}
