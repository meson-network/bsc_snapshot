package http_middleware

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/meson-network/bsc-data-file-utils/basic"
	"github.com/meson-network/bsc-data-file-utils/plugin/sqldb_plugin"
	"github.com/meson-network/bsc-data-file-utils/src/common/http"
	"github.com/meson-network/bsc-data-file-utils/src/common/limiter"
	"github.com/meson-network/bsc-data-file-utils/src/common/token_mgr"
	"github.com/meson-network/bsc-data-file-utils/src/user_mgr"
)

// check the correctness of the token format
func MID_TokenPreCheck(super_token_required bool) func(next echo.HandlerFunc) echo.HandlerFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := http.GetBearToken(c.Request().Header)
			if token == "" {
				return errors.New("token required")
			}

			if token_mgr.TokenMgr == nil {
				return errors.New("token error caused by system program error")
			}

			// check token format
			is_token, is_super_token := token_mgr.TokenMgr.CheckToken(token)
			if !is_token {
				basic.Logger.Debugln("token format error")
				return errors.New("token error")
			}

			if super_token_required && !is_super_token {
				basic.Logger.Debugln("super token error")
				return errors.New("super token error")
			}

			c.Set("token", token)
			c.Set("is_super_token", is_super_token)

			// no error continue next middleware call
			if err := next(c); err != nil {
				c.Error(err)
			}

			return nil
		}
	}
}

// MID_TokenPreCheck done required
func MID_TokenUser() func(next echo.HandlerFunc) echo.HandlerFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// check token
			token_i := c.Get("token")
			if token_i == nil {
				return errors.New("MID_TokenPreCheck must be used before MID_TokenUser ")
			}
			token := token_i.(string)

			limit := 1
			offset := 0

			userResult, err := user_mgr.QueryUser(sqldb_plugin.GetInstance(), nil, &token, nil, nil, nil, &limit, &offset, true, true)
			if err == nil && len(userResult.Users) > 0 {
				if userResult.Users[0].Forbidden {
					return errors.New("user forbidden")
				}
				c.Set("userInfo", userResult.Users[0])
			} else {
				return errors.New("token user not exist")
			}

			// no error continue next middleware call
			if err := next(c); err != nil {
				c.Error(err)
			}

			return nil
		}
	}

}

// MID_TokenUser done required
func MID_HasAnyRole(roles []string) func(echo.HandlerFunc) echo.HandlerFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if len(roles) == 0 {
				// bypass
			} else {
				user := GetUserInfo(c)
				if user == nil {
					return errors.New("user not exist")
				}

				if !user.HasOneOfRoles(roles) {
					return errors.New("has no correct role")
				}
			}

			// no error continue next middleware call
			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}

}

// MID_TokenUser done required
func MID_HasAllRoles(roles []string) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if len(roles) == 0 {
				// bypass
			} else {
				user := GetUserInfo(c)
				if user == nil {
					return errors.New("user not exist")
				}
				if !user.HasAllRoles(roles) {
					return errors.New("has no correct role")
				}
			}

			// no error continue next middleware call
			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}

}

// MID_TokenUser done required
func MID_HasAnyPermission(permissions []string) func(echo.HandlerFunc) echo.HandlerFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if len(permissions) == 0 {
				// bypass
			} else {
				user := GetUserInfo(c)
				if user == nil {
					return errors.New("user not exist")
				}
				if !user.HasOneOfPermissions(permissions) {
					return errors.New("has no correct permission")
				}
			}

			// no error continue next middleware call
			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}
}

// MID_TokenUser done required
func MID_HasAllPermissions(permissions []string) func(echo.HandlerFunc) echo.HandlerFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if len(permissions) == 0 {
				// bypass
			} else {
				user := GetUserInfo(c)
				if user == nil {
					return errors.New("user not exist")
				}
				if !user.HasAllPermissions(permissions) {
					return errors.New("has no correct permission")
				}
			}

			// no error continue next middleware call
			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}

}

// /////////////////////////////////////////////////////////////////////////

const default_speed_limiter_duration_second = 300
const default_speed_limiter_count = 300

// /////////////////////////////////////////////////////////////////////////
// MID_IP_Action speed limiter
func MID_IP_Action_SL(duration_second int64, Count int) func(echo.HandlerFunc) echo.HandlerFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			limited_action := c.Path() + c.RealIP()
			if !limiter.Allow(limited_action, duration_second, Count) {
				return errors.New("ip action overlimit:" + strconv.Itoa(Count) + "/" + strconv.Itoa(int(duration_second)))
			}

			// no error continue next middleware call
			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}
}

// MID_Default_IP_Action speed limiter
func MID_Default_IP_Action_SL() func(echo.HandlerFunc) echo.HandlerFunc {
	return MID_IP_Action_SL(default_speed_limiter_duration_second, default_speed_limiter_count)
}

// ///////////////////////////////////////////////////////////////////////
// MID_Token_IP_Action speed limiter
func MID_Token_IP_Action_SL(duration_second int64, Count int, bypass_super bool) func(echo.HandlerFunc) echo.HandlerFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// check token
			token_i := c.Get("token")
			is_super_token_i := c.Get("is_super_token")

			if token_i == nil || is_super_token_i == nil {
				return errors.New("MID_TokenPreCheck must be used before MID_Token_IP_SpeedLimiter ")
			}

			token := token_i.(string)
			is_super_token := is_super_token_i.(bool)

			if !is_super_token || !bypass_super {
				// token limiter
				limited_token_action := c.Path() + ":" + token
				if !limiter.Allow(limited_token_action, duration_second, Count) {
					return errors.New("token_ip token action overlimit:" + strconv.Itoa(Count) + "/" + strconv.Itoa(int(duration_second)))
				}

				// ip limiter
				limited_ip_action := c.Path() + ":" + c.RealIP()
				if !limiter.Allow(limited_ip_action, duration_second, Count) {
					return errors.New("token_ip ip action overlimit:" + strconv.Itoa(Count) + "/" + strconv.Itoa(int(duration_second)))
				}
			}

			// no error continue next middleware call
			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}
}

// MID_Default_Token_IP speed limiter
func MID_Default_Token_IP_SL() func(echo.HandlerFunc) echo.HandlerFunc {
	return MID_Token_IP_Action_SL(default_speed_limiter_duration_second, default_speed_limiter_count, true)
}

// ///////////////////////////////////////////////////////////////////////

func GetUserInfo(c echo.Context) *user_mgr.UserModel {
	v := c.Get("userInfo")
	if v == nil {
		return nil
	}
	return v.(*user_mgr.UserModel)
}
