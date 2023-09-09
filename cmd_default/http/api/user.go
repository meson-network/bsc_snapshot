package api

import (
	"net/http"
	"strings"

	"github.com/coreservice-io/utils/hash_util"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/basic/config"
	"github.com/meson-network/bsc-data-file-utils/cmd_default/http/api/http_middleware"
	"github.com/meson-network/bsc-data-file-utils/plugin/geo_ip_plugin"
	"github.com/meson-network/bsc-data-file-utils/plugin/mail_plugin"
	"github.com/meson-network/bsc-data-file-utils/plugin/sqldb_plugin"
	"github.com/meson-network/bsc-data-file-utils/src/common/captcha"
	"github.com/meson-network/bsc-data-file-utils/src/common/http/api"
	"github.com/meson-network/bsc-data-file-utils/src/common/json"
	"github.com/meson-network/bsc-data-file-utils/src/common/limiter"
	"github.com/meson-network/bsc-data-file-utils/src/common/token_mgr"
	"github.com/meson-network/bsc-data-file-utils/src/common/validator"
	"github.com/meson-network/bsc-data-file-utils/src/common/vcode"
	"github.com/meson-network/bsc-data-file-utils/src/user_mgr"
	"gorm.io/gorm"
)

type User struct {
	Id                      int64    `json:"id"`
	Email                   string   `json:"email"`
	Password                string   `json:"password"`
	Token                   string   `json:"token"`
	Is_super_token          bool     `json:"is_super_token"`
	Forbidden               bool     `json:"forbidden"`
	Roles                   []string `json:"roles"`
	Permissions             []string `json:"permissions"`
	Register_ip             string   `json:"register_ip"`
	Register_country_code   string   `json:"register_country_code"`
	Register_continent_code string   `json:"register_continent_code"`
	Register_region         string   `json:"register_region"`
	Created_unixtime        int64    `json:"created_unixtime"`
	Update_unixtime         int64    `json:"update_unixtime"`
}

type Msg_Req_EmailVCode struct {
	Email      string `json:"email"`      // required
	Vcode_len  int    `json:"vcode_len"`  // required
	Captcha_id string `json:"captcha_id"` // required
	Captcha    string `json:"captcha"`    // required
}

// @Msg_Resp_Captcha
type Msg_Resp_Captcha struct {
	api.API_META_STATUS
	Id      string `json:"id"`
	Content string `json:"content"`
}

// @Msg_Resp_User_Auth_Config
type Msg_Resp_User_Auth_Config struct {
	api.API_META_STATUS
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

// @Description Msg_Req_UserLogin
type Msg_Req_UserLogin struct {
	Email      string `json:"email"`      // required
	Password   string `json:"password"`   // required
	Captcha_id string `json:"captcha_id"` // required
	Captcha    string `json:"captcha"`    // required
}

// @Description Msg_Resp_Token
type Msg_Resp_Token struct {
	api.API_META_STATUS
	Token string `json:"token"`
}

// @Description Msg_Resp_UserInfo
type Msg_Resp_UserInfo struct {
	api.API_META_STATUS
	User *User `json:"user"`
}

// @Description Msg_Req_UserResetPassword
type Msg_Req_UserResetPassword struct {
	Email    string `json:"email"`    // required
	Password string `json:"password"` // required
	Vcode    string `json:"vcode"`    // required
}

// @Description Msg_Req_QueryNodeUser_Filter
type Msg_Req_QueryUser_Filter struct {
	Id            *int64  `json:"id"`            // optional
	Email_pattern *string `json:"email_pattern"` // optional
	Token         *string `json:"token"`         // optional
	Forbidden     *bool   `json:"forbidden"`     // optional
}

// @Description Msg_Req_QueryNodeUser
type Msg_Req_QueryUser struct {
	Filter Msg_Req_QueryUser_Filter `json:"filter"` // required
	Limit  *int                     `json:"limit"`  // optional
	Offset *int                     `json:"offset"` // optional
}

// @Description Msg_Resp_QueryNodeUser
type Msg_Resp_QueryUser struct {
	api.API_META_STATUS
	Data        []*User `json:"data"`
	Total_count int64   `json:"total_count"`
}

// @Description Msg_Req_UpdateNodeUser_Filter
type Msg_Req_UpdateUser_Filter struct {
	Id []int64 `json:"id"` // required
}

// @Description Msg_Req_UpdateNodeUser_To
type Msg_Req_UpdateUser_To struct {
	Forbidden   *bool     `json:"forbidden"`   // optional
	Roles       *[]string `json:"roles"`       // optional
	Permissions *[]string `json:"permissions"` // optional
}

// @Description Msg_Req_UpdateUser
type Msg_Req_UpdateUser struct {
	Filter Msg_Req_UpdateUser_Filter `json:"filter"` // required
	Update Msg_Req_UpdateUser_To     `json:"update"` // required
}

// create
// @Description Msg_Req_RegisterNodeUser
type Msg_Req_RegisterUser struct {
	Email    string `json:"email"`    // required
	Password string `json:"password"` // required
	Vcode    string `json:"vcode"`    // required
}

// create
// @Description Msg_Req_RegisterNodeUser
type Msg_Req_CreateUser struct {
	Email       string   `json:"email"`       // required
	Password    string   `json:"password"`    // required
	Roles       []string `json:"roles"`       // required
	Permissions []string `json:"permissions"` // required
}

// @Description Msg_Req_CheckEmail
type Msg_Req_CheckEmail struct {
	Email string `json:"email"` // required
}

func configUser(httpServer *echo.Echo) {
	//general
	httpServer.GET("/api/user/captcha", getCaptchaHandler, http_middleware.MID_Default_IP_Action_SL())
	httpServer.POST("/api/user/email_vcode", userEmailVCodeHandler, http_middleware.MID_Default_IP_Action_SL())

	//config
	httpServer.GET("/api/user/auth_config", authConfigHandler, http_middleware.MID_Default_IP_Action_SL())
	// user
	httpServer.POST("/api/user/login", userLoginHandler, http_middleware.MID_Default_IP_Action_SL())
	httpServer.POST("/api/user/register", userRegisterHandler, http_middleware.MID_Default_IP_Action_SL())
	httpServer.POST("/api/user/reset_password", userResetPasswordHandler, http_middleware.MID_Default_IP_Action_SL())
	httpServer.GET("/api/user/info", userInfoHandler, http_middleware.MID_TokenPreCheck(false), http_middleware.MID_Default_Token_IP_SL(), http_middleware.MID_TokenUser())

	// admin|readonly
	httpServer.POST("/api/user/admin/query", userQueryHandler, http_middleware.MID_TokenPreCheck(true),
		http_middleware.MID_Default_Token_IP_SL(), http_middleware.MID_TokenUser(),
		http_middleware.MID_HasAnyRole([]string{user_mgr.USER_ROLE_ADMIN, user_mgr.USER_ROLE_READONLY}))

	httpServer.POST("/api/user/admin/update", userUpdateHandler, http_middleware.MID_TokenPreCheck(true),
		http_middleware.MID_Default_Token_IP_SL(), http_middleware.MID_TokenUser(),
		http_middleware.MID_HasAllRoles([]string{user_mgr.USER_ROLE_ADMIN}))

	httpServer.POST("/api/user/admin/create", userCreateHandler, http_middleware.MID_TokenPreCheck(true),
		http_middleware.MID_Default_Token_IP_SL(), http_middleware.MID_TokenUser(),
		http_middleware.MID_HasAllRoles([]string{user_mgr.USER_ROLE_ADMIN}))

}

// @Summary      get captcha
// @Tags         user
// @Produce      json
// @response 	 200 {object} Msg_Resp_Captcha "result"
// @Router       /api/user/captcha [get]
func getCaptchaHandler(ctx echo.Context) error {

	res := &Msg_Resp_Captcha{}

	id, base64Code, err := captcha.GenCaptcha()
	if err != nil {
		// error gen captcha
		res.MetaStatus(-1, "gen captcha err:"+err.Error())
		return ctx.JSON(http.StatusOK, res)
	}
	if id == "" || base64Code == "" {
		// error gen captcha
		res.MetaStatus(-1, "gen captcha err")
		return ctx.JSON(http.StatusOK, res)
	}

	res.MetaStatus(1, "success")
	res.Id = id
	res.Content = base64Code
	return ctx.JSON(http.StatusOK, res)
}

// @Summary      send email VCode
// @Tags         user
// @Accept       json
// @Param        msg  body  Msg_Req_EmailVCode  true  "get email vcode"
// @Produce      json
// @Success      200 {object} api.API_META_STATUS "result"
// @Router       /api/user/email_vcode [post]
func userEmailVCodeHandler(ctx echo.Context) error {
	var msg Msg_Req_EmailVCode
	res := &api.API_META_STATUS{}
	if err := ctx.Bind(&msg); err != nil {
		res.MetaStatus(-1, "post data error")
		return ctx.JSON(http.StatusOK, res)
	}

	// captcha
	if !captcha.VerifyCaptcha(msg.Captcha_id, msg.Captcha) {
		res.MetaStatus(-2, "captcha error")
		return ctx.JSON(http.StatusOK, res)
	}

	// validate email
	email_err := validator.ValidateEmail(msg.Email)
	if email_err != nil {
		res.MetaStatus(-2, "email format error")
		return ctx.JSON(http.StatusOK, res)
	}

	// send vcode
	if !limiter.Allow(msg.Email, 20, 1) {
		res.MetaStatus(-3, "too frequent")
		return ctx.JSON(http.StatusOK, res)
	}

	code, err := vcode.GenVCode(msg.Email, msg.Vcode_len, 600)
	if err != nil {
		res.MetaStatus(-4, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}

	// send
	fromText := config.Get_config().Toml_config.Smtp.From_email
	toAddress := msg.Email
	subject := "Your verification code 'vcode', it will expire in 10 minutes "
	body := "" + code
	err = mail_plugin.GetInstance().Send(fromText, toAddress, subject, body)
	if err != nil {
		basic.Logger.Errorln("send VCode error:", err)
		res.MetaStatus(-1, "send VCode failed")
		return ctx.JSON(http.StatusOK, res)
	}

	res.MetaStatus(1, "success")
	return ctx.JSON(http.StatusOK, res)
}

// @Summary      get role setting
// @Tags         user
// @Produce      json
// @Success      200 {object} Msg_Resp_User_Auth_Config "result"
// @Router       /api/user/auth_config [get]
func authConfigHandler(ctx echo.Context) error {
	res := &Msg_Resp_User_Auth_Config{}
	res.MetaStatus(1, "success")
	res.Roles = user_mgr.UserRoles
	res.Permissions = user_mgr.UserPermissions
	return ctx.JSON(http.StatusOK, res)
}

// @Summary      get user info
// @Tags         user
// @Security 	ApiKeyAuth
// @Produce      json
// @Success      200 {object} Msg_Resp_UserInfo "result"
// @Router       /api/user/info [get]
func userInfoHandler(ctx echo.Context) error {
	res := &Msg_Resp_UserInfo{}

	userInfo := http_middleware.GetUserInfo(ctx)
	if userInfo == nil {
		res.MetaStatus(-1, "user error")
		return ctx.JSON(http.StatusOK, res)
	}

	res.User = &User{}
	err := copier.Copy(&res.User, &userInfo)
	if err != nil {
		// should never happen
		basic.Logger.Errorln("userInfoHandler copier.Copy error", err)
		res.MetaStatus(-1, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}
	res.User.Password = ""
	var roles []string
	err = json.Unmarshal([]byte(userInfo.Roles), &roles)
	if err != nil {
		res.MetaStatus(-2, "json.Unmarshal user roles error")
		return ctx.JSON(http.StatusOK, res)
	}
	var permissions []string
	err = json.Unmarshal([]byte(userInfo.Permissions), &permissions)
	if err != nil {
		res.MetaStatus(-3, "json.Unmarshal user permissions error")
		return ctx.JSON(http.StatusOK, res)
	}

	res.User.Roles = roles
	res.User.Permissions = permissions
	res.User.Password = ""

	// ip
	ipInfo, err := geo_ip_plugin.GetInstance().GetInfo(userInfo.Register_ipv4)
	if err == nil {
		res.User.Register_continent_code = ipInfo.Continent_code
		res.User.Register_country_code = ipInfo.Country_code
		res.User.Register_region = ipInfo.Region
	} else {
		res.User.Register_continent_code = "ZZ"
		res.User.Register_country_code = "ZZ"
		res.User.Register_region = "ZZ"
	}

	res.MetaStatus(1, "success")
	return ctx.JSON(http.StatusOK, res)
}

// @Summary      user login
// @Tags         user
// @Accept       json
// @Param        msg  body  Msg_Req_UserLogin  true  "user login"
// @Produce      json
// @Success      200 {object} Msg_Resp_UserInfo "result"
// @Router       /api/user/login [post]
func userLoginHandler(ctx echo.Context) error {
	var msg Msg_Req_UserLogin
	res := &Msg_Resp_Token{}
	if err := ctx.Bind(&msg); err != nil {
		res.MetaStatus(-1, "post data error")
		return ctx.JSON(http.StatusOK, res)
	}

	// captcha
	if !captcha.VerifyCaptcha(msg.Captcha_id, msg.Captcha) {
		res.MetaStatus(-2, "captcha error")
		return ctx.JSON(http.StatusOK, res)
	}

	// validate email
	verr := validator.ValidateEmail(msg.Email)
	if verr != nil {
		res.MetaStatus(-3, "email format error")
		return ctx.JSON(http.StatusOK, res)
	}

	/////

	limit := 1
	offset := 0
	// find user
	userResult, err := user_mgr.QueryUser(sqldb_plugin.GetInstance(), nil, nil, nil, &msg.Email, nil, &limit, &offset, true, true)
	if err != nil {
		res.MetaStatus(-4, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}
	if len(userResult.Users) == 0 {
		res.MetaStatus(-5, "login email/password error")
		return ctx.JSON(http.StatusOK, res)
	}
	userInfo := userResult.Users[0]
	if userInfo.Forbidden {
		res.MetaStatus(-6, "user forbidden")
		return ctx.JSON(http.StatusOK, res)
	}

	// password
	if hash_util.SHA256String(msg.Password) != userInfo.Password {
		res.MetaStatus(-7, "login email/password error")
		return ctx.JSON(http.StatusOK, res)
	}

	res.Token = userInfo.Token
	res.MetaStatus(1, "success")
	return ctx.JSON(http.StatusOK, res)
}

// @Summary      user register
// @Description  user register
// @Tags         user
// @Accept       json
// @Param        msg  body  Msg_Req_RegisterUser true  "new user info"
// @Produce      json
// @Success      200 {object} Msg_Resp_Token "result"
// @Router       /api/user/register [post]
func userRegisterHandler(ctx echo.Context) error {
	var msg Msg_Req_RegisterUser
	res := &Msg_Resp_Token{}

	if err := ctx.Bind(&msg); err != nil {
		res.MetaStatus(-1, "bad request data")
		return ctx.JSON(http.StatusOK, res)
	}

	// validate email
	email_err := validator.ValidateEmail(msg.Email)
	if email_err != nil {
		res.MetaStatus(-3, "email format error")
		return ctx.JSON(http.StatusOK, res)
	}

	// validate password
	passwd_err := validator.ValidatePassword(msg.Password)
	if passwd_err != nil {
		res.MetaStatus(-4, "password format error:must contain number and letter, special character is optional,length 6-20")
		return ctx.JSON(http.StatusOK, res)
	}

	// check vcode
	vcode_ := strings.TrimSpace(msg.Vcode)
	if !vcode.ValidateVCodeWithLen(msg.Email, vcode_, 16) {
		res.MetaStatus(-5, "vcode error, wrong or expires")
		return ctx.JSON(http.StatusOK, res)
	}

	// find user
	limit := 1
	offset := 0

	userResult, err := user_mgr.QueryUser(sqldb_plugin.GetInstance(), nil, nil, nil, &msg.Email, nil, &limit, &offset, false, true)
	if err != nil {
		res.MetaStatus(-6, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}
	if len(userResult.Users) > 0 {
		res.MetaStatus(-7, "user already exist")
		return ctx.JSON(http.StatusOK, res)
	}

	// ip
	remoteIp := ctx.RealIP()

	var newUser *user_mgr.UserModel
	err = sqldb_plugin.GetInstance().Transaction(func(tx *gorm.DB) error {
		newUser, err = user_mgr.CreateUser(tx, msg.Email, msg.Password, false, []string{user_mgr.USER_ROLE_USER}, []string{}, remoteIp)
		if err != nil {
			return err
		}

		// other action

		return nil
	})
	if err != nil {
		res.MetaStatus(-10, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}

	//refresh cache
	user_mgr.QueryUser(sqldb_plugin.GetInstance(), nil, nil, nil, &msg.Email, nil, &limit, &offset, true, true)

	res.Token = newUser.Token
	res.MetaStatus(1, "success")
	return ctx.JSON(http.StatusOK, res)
}

// @Summary      user reset password
// @Tags         user
// @Accept       json
// @Param        msg  body  Msg_Req_UserResetPassword  true  "user reset password"
// @Produce      json
// @Success      200 {object} api.API_META_STATUS "result"
// @Router       /api/user/reset_password [post]
func userResetPasswordHandler(ctx echo.Context) error {
	var msg Msg_Req_UserResetPassword
	res := &api.API_META_STATUS{}
	if err := ctx.Bind(&msg); err != nil {
		res.MetaStatus(-1, "post data error")
		return ctx.JSON(http.StatusOK, res)
	}

	// validate email
	email_err := validator.ValidateEmail(msg.Email)
	if email_err != nil {
		res.MetaStatus(-2, "email format error")
		return ctx.JSON(http.StatusOK, res)
	}

	// validate password
	pass_err := validator.ValidatePassword(msg.Password)
	if pass_err != nil {
		res.MetaStatus(-3, "password format error")
		return ctx.JSON(http.StatusOK, res)
	}

	// check vcode
	vcode_ := strings.TrimSpace(msg.Vcode)
	if !vcode.ValidateVCodeWithLen(msg.Email, vcode_, 16) {
		res.MetaStatus(-4, "vcode error, wrong or expire")
		return ctx.JSON(http.StatusOK, res)
	}

	limit := 1
	offset := 0

	userResult, err := user_mgr.QueryUser(sqldb_plugin.GetInstance(), nil, nil, nil, &msg.Email, nil, &limit, &offset, false, false)
	if err != nil {
		res.MetaStatus(-5, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}
	if len(userResult.Users) == 0 {
		res.MetaStatus(-6, "user not exist")
		return ctx.JSON(http.StatusOK, res)
	}
	user := userResult.Users[0]

	if user.Password == hash_util.SHA256String(msg.Password) {
		res.MetaStatus(-7, "same as old password")
		return ctx.JSON(http.StatusOK, res)
	}

	err = user_mgr.UpdateUser(sqldb_plugin.GetInstance(), map[string]interface{}{
		"password": hash_util.SHA256String(msg.Password),
	}, user.Id)
	if err != nil {
		res.MetaStatus(-8, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}

	user_mgr.QueryUser(sqldb_plugin.GetInstance(), nil, nil, nil, &msg.Email, nil, &limit, &offset, false, true)

	res.MetaStatus(1, "success")
	return ctx.JSON(http.StatusOK, res)
}

// @Summary      query user
// @Tags         user
// @Security ApiKeyAuth
// @Accept       json
// @Param        msg  body  Msg_Req_QueryUser  true  "query user condition"
// @Produce      json
// @Success      200 {object} Msg_Resp_QueryUser "result"
// @Router       /api/user/admin/query [post]
func userQueryHandler(ctx echo.Context) error {
	var msg Msg_Req_QueryUser
	res := &Msg_Resp_QueryUser{}
	res.Data = []*User{}

	if err := ctx.Bind(&msg); err != nil {
		res.MetaStatus(-1, "post data error")
		return ctx.JSON(http.StatusOK, res)
	}

	userResult, err := user_mgr.QueryUser(sqldb_plugin.GetInstance(), msg.Filter.Id, msg.Filter.Token, msg.Filter.Email_pattern, nil, msg.Filter.Forbidden, msg.Limit, msg.Offset, false, false)
	if err != nil {
		res.MetaStatus(-4, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}

	userInfo := http_middleware.GetUserInfo(ctx)

	for _, um := range userResult.Users {
		um.Password = ""

		msguser := User{}
		err = copier.Copy(&msguser, &um)
		if err != nil {
			// should never happen
			basic.Logger.Errorln("userQueryHandler copier.Copy error", err)
			res.MetaStatus(-1, err.Error())
			return ctx.JSON(http.StatusOK, res)
		}

		// get location
		if um.Register_ipv4 != "" {
			ipInfo, err := geo_ip_plugin.GetInstance().GetInfo(um.Register_ipv4)
			if err != nil {
				basic.Logger.Errorln("user register ip can't be reserved. ip:" + um.Register_ipv4)
			} else {
				msguser.Register_country_code = ipInfo.Country_code
				msguser.Register_continent_code = ipInfo.Continent_code
				msguser.Register_region = ipInfo.Region
			}
		}

		var roles []string
		err = json.Unmarshal([]byte(um.Roles), &roles)
		if err != nil {
			res.MetaStatus(-1, "json.Unmarshal user roles error")
			return ctx.JSON(http.StatusOK, res)
		}
		var permissions []string
		err = json.Unmarshal([]byte(um.Permissions), &permissions)
		if err != nil {
			res.MetaStatus(-1, "json.Unmarshal user permissions error")
			return ctx.JSON(http.StatusOK, res)
		}

		msguser.Roles = roles
		msguser.Permissions = permissions

		_, isSuperToken := token_mgr.TokenMgr.CheckToken(um.Token)
		msguser.Is_super_token = isSuperToken

		//for safety reason token only visable to self
		if msguser.Id != userInfo.Id {
			msguser.Token = ""
		}

		res.Data = append(res.Data, &msguser)
	}

	res.Total_count = userResult.Total_count
	res.MetaStatus(1, "success")
	return ctx.JSON(http.StatusOK, res)
}

// @Summary      admin update user
// @Tags         user
// @Security ApiKeyAuth
// @Accept		 json
// @Param        msg  body  Msg_Req_UpdateUser true  "update"
// @Produce      json
// @Success      200 {object} api.API_META_STATUS "result"
// @Router       /api/user/admin/update [post]
func userUpdateHandler(ctx echo.Context) error {
	var msg Msg_Req_UpdateUser
	res := &api.API_META_STATUS{}

	if err := ctx.Bind(&msg); err != nil {
		res.MetaStatus(-1, "post data error")
		return ctx.JSON(http.StatusOK, res)
	}

	if len(msg.Filter.Id) == 0 {
		res.MetaStatus(-2, "user id needed")
		return ctx.JSON(http.StatusOK, res)
	}
	if len(msg.Filter.Id) > 1 {
		res.MetaStatus(-3, "only support one user id")
		return ctx.JSON(http.StatusOK, res)
	}

	userInfo := http_middleware.GetUserInfo(ctx)
	id := msg.Filter.Id[0]
	if msg.Update.Forbidden != nil && *msg.Update.Forbidden && userInfo.Id == id {
		res.MetaStatus(-6, "can not forbidden self")
		return ctx.JSON(http.StatusOK, res)
	}

	limit := 1
	offset := 0

	result, err := user_mgr.QueryUser(sqldb_plugin.GetInstance(), &msg.Filter.Id[0], nil, nil, nil, nil, &limit, &offset, true, true)
	if err != nil {
		res.MetaStatus(-7, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}
	if len(result.Users) == 0 {
		res.MetaStatus(-8, "user not exist")
		return ctx.JSON(http.StatusOK, res)
	}

	update_map := map[string]interface{}{}
	if msg.Update.Forbidden != nil {
		update_map["forbidden"] = *msg.Update.Forbidden
	}

	if msg.Update.Roles != nil {
		if !user_mgr.RolesDefined(*msg.Update.Roles) {
			res.MetaStatus(-9, "roles error")
			return ctx.JSON(http.StatusOK, res)
		}
		roles_json, err := json.Marshal(msg.Update.Roles)
		if err != nil {
			res.MetaStatus(-10, "roles marshal error")
			return ctx.JSON(http.StatusOK, res)
		}
		update_map["roles"] = string(roles_json)
	}

	if msg.Update.Permissions != nil {
		if !user_mgr.PermissionsDefined(*msg.Update.Permissions) {
			res.MetaStatus(-11, "permissions error")
			return ctx.JSON(http.StatusOK, res)
		}
		permissions_json, err := json.Marshal(msg.Update.Permissions)
		if err != nil {
			res.MetaStatus(-12, "permissions marshal error")
			return ctx.JSON(http.StatusOK, res)
		}
		update_map["permissions"] = string(permissions_json)
	}

	err = user_mgr.UpdateUser(sqldb_plugin.GetInstance(), update_map, id)
	if err != nil {
		res.MetaStatus(-13, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}

	// cache
	user_mgr.QueryUser(sqldb_plugin.GetInstance(), nil, &result.Users[0].Token, nil, nil, nil, &limit, &offset, false, true)

	res.MetaStatus(1, "success")
	return ctx.JSON(http.StatusOK, res)
}

// @Summary      admin create user
// @Tags         user
// @Security ApiKeyAuth
// @Accept		 json
// @Param        msg  body  Msg_Req_CreateUser true  "update"
// @Produce      json
// @Success      200 {object} api.API_META_STATUS "result"
// @Router       /api/user/admin/create [post]
func userCreateHandler(ctx echo.Context) error {
	var msg Msg_Req_CreateUser
	res := &api.API_META_STATUS{}

	if err := ctx.Bind(&msg); err != nil {
		res.MetaStatus(-1, "post data error")
		return ctx.JSON(http.StatusOK, res)
	}

	// validate email
	email_err := validator.ValidateEmail(msg.Email)
	if email_err != nil {
		res.MetaStatus(-3, "email format error")
		return ctx.JSON(http.StatusOK, res)
	}

	// validate password
	passwd_err := validator.ValidatePassword(msg.Password)
	if passwd_err != nil {
		res.MetaStatus(-4, "password format error:must contain number and letter, special character is optional,length 6-20")
		return ctx.JSON(http.StatusOK, res)
	}

	////////
	limit := 1
	offset := 0

	// check user exist
	result, err := user_mgr.QueryUser(sqldb_plugin.GetInstance(), nil, nil, nil, &msg.Email, nil, &limit, &offset, true, true)
	if err != nil {
		res.MetaStatus(-7, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}
	if len(result.Users) > 0 {
		res.MetaStatus(-8, "user already exist")
		return ctx.JSON(http.StatusOK, res)
	}

	isSuperUser := false
	for _, v := range msg.Roles {
		if v == user_mgr.USER_ROLE_ADMIN || v == user_mgr.USER_ROLE_READONLY {
			isSuperUser = true
			break
		}
	}

	// ip
	remoteIp := ctx.RealIP()
	err = sqldb_plugin.GetInstance().Transaction(func(tx *gorm.DB) error {
		_, err = user_mgr.CreateUser(tx, msg.Email, msg.Password, isSuperUser, msg.Roles, msg.Permissions, remoteIp)
		if err != nil {
			return err
		}
		// other action
		return nil
	})
	if err != nil {
		res.MetaStatus(-10, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}

	//refresh cache
	user_mgr.QueryUser(sqldb_plugin.GetInstance(), nil, nil, nil, &msg.Email, nil, &limit, &offset, true, true)

	res.MetaStatus(1, "success")
	return ctx.JSON(http.StatusOK, res)
}
