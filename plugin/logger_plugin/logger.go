package logger_plugin

import (
	"fmt"

	"github.com/coreservice-io/log"
	"github.com/coreservice-io/logrus_log"
)

var instanceMap = map[string]log.Logger{}

func GetInstance() log.Logger {
	return instanceMap["default"]
}

func GetInstance_(name string) log.Logger {
	return instanceMap[name]
}

func Init(log_abs_path string) error {
	return Init_("default", log_abs_path)
}

func Init_(name string, log_abs_path string) error {

	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("email instance <%s> has already been initialized", name)
	}

	Logger, llerr := logrus_log.New(log_abs_path, 2, 20, 30)
	if llerr != nil {
		panic("Error:" + llerr.Error())
	}

	instanceMap[name] = Logger
	return nil
}
