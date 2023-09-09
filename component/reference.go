package component

import (
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/plugin/reference_plugin"
)

func InitReference() error {
	if err := reference_plugin.Init(); err == nil {
		basic.Logger.Infoln("### InitReference success")
		return nil
	} else {
		return err
	}
}
