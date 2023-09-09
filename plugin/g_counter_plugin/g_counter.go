package g_counter_plugin

import (
	"fmt"

	general_counter "github.com/coreservice-io/general-counter"
	"github.com/coreservice-io/log"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"gorm.io/gorm"
)

var instanceMap = map[string]*general_counter.GeneralCounter{}

func GetInstance() *general_counter.GeneralCounter {
	return GetInstance_("default")
}

func GetInstance_(name string) *general_counter.GeneralCounter {
	gcounter := instanceMap[name]
	if gcounter == nil {
		basic.Logger.Errorln(name + " general counter plugin null")
	}
	return gcounter
}

func Init(gConfig *general_counter.GeneralCounterConfig, logger log.Logger) error {
	return Init_("default", gConfig, logger)
}

func DBMigrate(db *gorm.DB) {
	db.AutoMigrate(&general_counter.GCounterModel{},
		&general_counter.GCounterDailyAggModel{},
		&general_counter.GCounterDetailModel{})
}

// Init a new instance.
//
//	If only need one instance, use empty name "". Use GetDefaultInstance() to get.
//	If you need several instance, run Init() with different <name>. Use GetInstance(<name>) to get.
func Init_(name string, gConfig *general_counter.GeneralCounterConfig, logger log.Logger) error {
	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("spr instance <%s> has already been initialized", name)
	}

	// ////// ini g_counter //////////////////////
	gcounter, err := general_counter.NewGeneralCounter(gConfig, logger)

	if err != nil {
		panic("init general_counter err:" + err.Error())
	}

	instanceMap[name] = gcounter

	return nil
}
