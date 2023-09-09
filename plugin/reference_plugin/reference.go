package reference_plugin

import (
	"fmt"

	"github.com/coreservice-io/reference"
)

var instanceMap = map[string]*reference.Reference{}

func GetInstance() *reference.Reference {
	return instanceMap["default"]
}

func GetInstance_(name string) *reference.Reference {
	return instanceMap[name]
}

func Init() error {
	return Init_("default")
}

// Init a new instance.
//  If only need one instance, use empty name "". Use GetDefaultInstance() to get.
//  If you need several instance, run Init() with different <name>. Use GetInstance(<name>) to get.
func Init_(name string) error {
	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("reference instance <%s> has already been initialized", name)
	}
	instanceMap[name] = reference.New()
	return nil
}
