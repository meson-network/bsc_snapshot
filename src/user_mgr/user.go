package user_mgr

import (
	"errors"

	"github.com/coreservice-io/utils/hash_util"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/plugin/redis_plugin"
	"github.com/meson-network/bsc-data-file-utils/src/common/json"
	"github.com/meson-network/bsc-data-file-utils/src/common/smart_cache"
	"github.com/meson-network/bsc-data-file-utils/src/common/token_mgr"
	"gorm.io/gorm"
)

func RolesToStr(roles ...string) string {
	r_str, _ := json.Marshal(roles)
	return string(r_str)
}

func PermissionsToStr(permissions ...string) string {
	p_str, _ := json.Marshal(permissions)
	return string(p_str)
}

func GenRandUserToken(isSuperUser bool) string {
	if isSuperUser {
		return token_mgr.TokenMgr.GenSuperToken()
	} else {
		return token_mgr.TokenMgr.GenToken()
	}
}

func CreateUser(tx *gorm.DB, email string, passwd string, isSuperUser bool, roles []string, permissions []string, ipv4 string) (*UserModel, error) {
	sha256_passwd := hash_util.SHA256String(passwd)
	token := GenRandUserToken(isSuperUser)

	// check roles
	if !RolesDefined(roles) {
		return nil, errors.New("roles not defined")
	}

	// check permissions
	if !PermissionsDefined(permissions) {
		return nil, errors.New("permissions not defined")
	}

	user := &UserModel{
		Email:           email,
		Password:        sha256_passwd,
		Token:           token,
		Roles:           RolesToStr(roles...),
		Permissions:     PermissionsToStr(permissions...),
		Roles_map:       make(map[string]string),
		Permissions_map: make(map[string]string),
		Forbidden:       false,
		Register_ipv4:   ipv4,
	}
	for _, role := range roles {
		user.Roles_map[role] = role
	}

	for _, p := range permissions {
		user.Permissions_map[p] = p
	}

	if err := tx.Table(TABLE_NAME_USER).Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(tx *gorm.DB, updateData map[string]interface{}, id int64) error {

	limit := 1
	offset := 0

	queryResult, err := QueryUser(tx, &id, nil, nil, nil, nil, &limit, &offset, false, false)
	if err != nil {
		basic.Logger.Errorln("UpdateNodeUser queryUsers error:", err, "id:", id)
		return err
	}
	if len(queryResult.Users) == 0 {
		return errors.New("user not exist")
	}

	update_result := tx.Table(TABLE_NAME_USER).Where("id =?", id).Updates(updateData)
	if update_result.Error != nil {
		return update_result.Error
	}

	if update_result.RowsAffected == 0 {
		return errors.New("0 raw affected")
	}

	// update cache , for fast api middleware token auth
	QueryUser(tx, nil, &queryResult.Users[0].Token, nil, nil, nil, &limit, &offset, false, true)

	return nil
}

type QueryUserResult struct {
	Users       []*UserModel
	Total_count int64
}

func QueryUser(tx *gorm.DB, id *int64, token *string, emailPattern *string, email *string, forbidden *bool, limit *int, offset *int, fromCache bool, updateCache bool) (*QueryUserResult, error) {

	if err := smart_cache.CheckLimitOffset(limit, offset); err != nil {
		return nil, err
	}

	if emailPattern != nil && email != nil {
		return &QueryUserResult{
			Users:       []*UserModel{},
			Total_count: 0,
		}, errors.New("emailPattern ,email :can't be set at the same time")
	}

	// gen_key
	ck := smart_cache.NewConnectKey("users")
	ck.C_Str_Ptr("token", token).
		C_Str_Ptr("emailPattern", emailPattern).
		C_Str_Ptr("email", email).
		C_Bool_Ptr("forbidden", forbidden).
		C_Int64_Ptr("id", id).
		C_Int_Ptr("limit", limit).
		C_Int_Ptr("offset", offset)

	key := redis_plugin.GetInstance().GenKey(ck.String())

	// ///
	resultHolderAlloc := func() *smart_cache.QueryResult {
		return &smart_cache.QueryResult{
			Result_holder: &QueryUserResult{
				Users:       []*UserModel{},
				Total_count: 0,
			},
			Found:   false,
			Has_err: false,
			Err_str: "",
		}
	}

	query := func(resultHolder *smart_cache.QueryResult) *smart_cache.QueryCacheTTL {

		queryResult := resultHolder.Result_holder.(*QueryUserResult)

		query := tx.Table(TABLE_NAME_USER)
		if id != nil {
			query.Where("id = ?", *id)
		}

		if token != nil {
			query.Where("token = ?", *token)
		}

		if emailPattern != nil {
			query.Where("email LIKE ?", "%"+*emailPattern+"%")
		}

		if email != nil {
			query.Where("email = ?", *email)
		}

		if forbidden != nil {
			query.Where("forbidden = ?", *forbidden)
		}

		query.Count(&queryResult.Total_count)
		if limit != nil {
			query.Limit(*limit)
		}
		if offset != nil {
			query.Offset(*offset)
		}

		err := query.Find(&queryResult.Users).Error
		if err != nil {
			resultHolder.Has_err = true
			resultHolder.Err_str = err.Error()
			resultHolder.Found = false
			return smart_cache.SlowQueryTTL_ERR
		}

		if queryResult.Total_count == 0 {
			return smart_cache.SlowQueryTTL_NOT_FOUND
		}

		// equip the related info
		for _, user := range queryResult.Users {

			user.Roles_map = make(map[string]string)
			var userRoles []string
			json.Unmarshal([]byte(user.Roles), &userRoles)
			for _, u_r := range userRoles {
				user.Roles_map[u_r] = u_r
			}

			user.Permissions_map = make(map[string]string)
			var userPermissions []string
			json.Unmarshal([]byte(user.Permissions), &userPermissions)
			for _, u_p := range userPermissions {
				user.Permissions_map[u_p] = u_p
			}
		}

		resultHolder.Found = true
		return smart_cache.SlowQueryTTL_Default

	}

	s_query := &smart_cache.SlowQuery{
		CacheTTL: smart_cache.SlowQueryTTL_Default,
		Query:    query,
	}

	//
	sq_result := smart_cache.SmartQueryCacheSlow(key, fromCache, updateCache, "QueryUser", resultHolderAlloc, s_query)
	if sq_result.Has_err {
		return nil, errors.New(sq_result.Err_str)
	} else {
		return sq_result.Result_holder.(*QueryUserResult), nil
	}

}
