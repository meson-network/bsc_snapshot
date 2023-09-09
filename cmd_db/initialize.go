package cmd_db

import (
	"errors"

	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/plugin/sqldb_plugin"
	"github.com/meson-network/bsc-data-file-utils/src/common/dbkv"
	"github.com/meson-network/bsc-data-file-utils/src/user_mgr"
	"gorm.io/gorm"
)

func Initialize() {
	StartDBComponent()

	//
	key := "db_initialized"
	err := sqldb_plugin.GetInstance().Transaction(func(tx *gorm.DB) error {

		value, err := dbkv.GetDBKV(tx, nil, &key, false, false)
		if err != nil {
			return errors.New("db error" + err.Error())
		}

		if value != nil {
			return errors.New("db already initialized")
		}

		// create your own data here which won't change in the future
		configAdmin(tx)

		// dbkv
		return dbkv.SetDBKV_Bool(tx, key, true, "db initialized mark")
	})

	if err != nil {
		basic.Logger.Errorln("db initialize error:", err)
		return
	} else {
		basic.Logger.Infoln("db initialized")
	}

}

// =====below data can be changed=====
var ini_admin_email = "admin@coreservice.com"
var ini_admin_password = "abc@123"
var ini_admin_roles = user_mgr.UserRoles
var ini_admin_permissions = user_mgr.UserPermissions

func configAdmin(tx *gorm.DB) {

	default_u_admin, err := user_mgr.CreateUser(tx, ini_admin_email, ini_admin_password, true, ini_admin_roles, ini_admin_permissions, "127.0.0.1")
	if err != nil {
		basic.Logger.Panicln(err)
	} else {
		basic.Logger.Infoln("default admin created:", default_u_admin)
	}
}
