package user_mgr

import (
	"github.com/meson-network/bsc-data-file-utils/src/common/data"
)

const (
	USER_ROLE_ADMIN    = "admin"
	USER_ROLE_USER     = "user"
	USER_ROLE_READONLY = "read_only"
)

var UserRoles = []string{USER_ROLE_ADMIN, USER_ROLE_USER, USER_ROLE_READONLY}
var UserPermissions = []string{}

// only return true if  permissions are all defined
func PermissionsDefined(permissions []string) bool {

	if len(permissions) == 0 {
		return true
	}

	for _, v := range permissions {
		if !data.InArray(v, UserPermissions) {
			return false
		}
	}
	return true
}

// only return true if  roles are all defined
func RolesDefined(roles []string) bool {

	if len(roles) == 0 {
		return true
	}

	for _, v := range roles {
		if !data.InArray(v, UserRoles) {
			return false
		}
	}
	return true
}

// user must have all the roles specified
func (u *UserModel) HasAllRoles(roles []string) bool {
	for _, role := range roles {
		if _, ok := u.Roles_map[role]; !ok {
			return false
		}
	}
	return true
}

// user must have all the permissions specified
func (u *UserModel) HasAllPermissions(permissions []string) bool {
	for _, p := range permissions {
		if _, ok := u.Permissions_map[p]; !ok {
			return false
		}
	}
	return true
}

func (u *UserModel) HasOneOfRoles(roles []string) bool {
	if len(roles) == 0 {
		return true
	}
	for _, role := range roles {
		if _, ok := u.Roles_map[role]; ok {
			return true
		}
	}

	return false
}

func (u *UserModel) HasOneOfPermissions(permissions []string) bool {
	if len(permissions) == 0 {
		return true
	}

	for _, p := range permissions {
		if _, ok := u.Permissions_map[p]; ok {
			return true
		}
	}

	return false
}
