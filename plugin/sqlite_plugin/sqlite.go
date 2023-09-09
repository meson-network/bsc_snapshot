package sqlite_plugin

import (
	"fmt"
	"time"

	"github.com/coreservice-io/gorm_log"
	"github.com/coreservice-io/log"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var instanceMap = map[string]*gorm.DB{}

func GetInstance() *gorm.DB {
	return GetInstance_("default")
}

func GetInstance_(name string) *gorm.DB {
	sqlite_i := instanceMap[name]
	if sqlite_i == nil {
		basic.Logger.Errorln(name + " sqlite plugin null")
	}
	return sqlite_i
}

type Config struct {
	Sqlite_abs_path string
}

func Init(dbConfig *Config, logger log.Logger) error {
	return Init_("default", dbConfig, logger)
}

// Init a new instance.
//
//	If only need one instance, use empty name "". Use GetDefaultInstance() to get.
//	If you need several instance, run Init() with different <name>. Use GetInstance(<name>) to get.
func Init_(name string, dbConfig *Config, logger log.Logger) error {
	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("sqlite instance <%s> has already been initialized", name)
	}

	//Level: Silent Error Warn Info. Info logs all record. Silent turns off log.
	db_log_level := gorm_log.Warn
	if logger.GetLevel() >= log.TraceLevel {
		db_log_level = gorm_log.Info
	}

	db, err := gorm.Open(sqlite.Open(dbConfig.Sqlite_abs_path), &gorm.Config{
		Logger: gorm_log.New_gormLocalLogger(logger, gorm_log.Config{
			SlowThreshold:             500 * time.Millisecond,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  db_log_level,
		}),
	})

	if err != nil {
		return err
	}

	_, err = db.DB()
	if err != nil {
		return err
	}

	instanceMap[name] = db

	return nil
}
