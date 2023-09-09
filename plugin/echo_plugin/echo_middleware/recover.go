package echo_middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type (
	// RecoverConfig defines the config for Recover middleware.
	RecoverConfig struct {
		OnPanic func(panic_err interface{})
	}
)

// RecoverWithConfig returns a Recover middleware with config.
// See: `Recover()`.
func RecoverWithConfig(config RecoverConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					if config.OnPanic != nil {
						config.OnPanic(r)
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
