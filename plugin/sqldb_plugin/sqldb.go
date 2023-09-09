package sqldb_plugin

import (
	"fmt"
	"strconv"
	"time"

	"github.com/coreservice-io/gorm_log"
	"github.com/coreservice-io/log"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var instanceMap = map[string]*gorm.DB{}

func GetInstance() *gorm.DB {
	return GetInstance_("default")
}

func GetInstance_(name string) *gorm.DB {
	db_i := instanceMap[name]
	if db_i == nil {
		basic.Logger.Errorln(name + " sqldb plugin null")
	}
	return db_i
}

/*
db_host
db_port
db_name
db_username
db_password
*/
type Config struct {
	Host     string
	Port     int
	DbName   string
	UserName string
	Password string
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
		return fmt.Errorf("sqldb instance <%s> has already been initialized", name)
	}

	//Level: Silent Error Warn Info. Info logs all record. Silent turns off log.
	db_log_level := gorm_log.Warn
	if logger.GetLevel() >= log.TraceLevel {
		db_log_level = gorm_log.Info
	}

	dsn := dbConfig.UserName + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port) + ")/" + dbConfig.DbName + "?charset=utf8mb4&loc=UTC"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gorm_log.New_gormLocalLogger(logger, gorm_log.Config{
			SlowThreshold:             500 * time.Millisecond,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  db_log_level,
		}),
	})

	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)

	instanceMap[name] = db

	return nil
}
