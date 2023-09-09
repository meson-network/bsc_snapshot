package component

import (
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/basic/config"
	"github.com/meson-network/bsc-data-file-utils/plugin/mail_plugin"
)

func InitSmtpMail(toml_conf *config.TomlConfig) error {

	if toml_conf.Smtp.Enable {

		mail_conf := mail_plugin.Config{
			Host:     toml_conf.Smtp.Host,
			Port:     toml_conf.Smtp.Port,
			Password: toml_conf.Smtp.Password,
			UserName: toml_conf.Smtp.Username,
		}

		basic.Logger.Infoln("init smtp mail plugin with config:", mail_conf)

		if err := mail_plugin.Init(&mail_conf); err == nil {
			basic.Logger.Infoln("### InitSmtpMail success")
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}
