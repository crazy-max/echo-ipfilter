package ipfilter

import (
	"fmt"
	"net"
	"net/http"

	"github.com/jpillora/ipfilter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Config defines the config for IPFilter middleware.
type Config struct {
	// Skipper defines a function to skip middleware.
	// default middleware.DefaultSkipper
	Skipper middleware.Skipper

	// WhiteList is an allowed ip list.
	WhiteList []string

	// BlackList is a disallowed ip list.
	BlackList []string

	// Block by default.
	BlockByDefault bool
}

// var DefaultConfig = IPFilterConfig{ is the default IPFilter middleware config config.
var DefaultConfig = Config{
	Skipper:        middleware.DefaultSkipper,
	BlockByDefault: false,
}

// Middleware returns an IPFilter middleware to
// filter requests by ip matching / blocking.
func Middleware() echo.MiddlewareFunc {
	return MiddlewareWithConfig(DefaultConfig)
}

// MiddlewareWithConfig returns an IPFilter middleware with config.
// See: `IPFilter()`.
func MiddlewareWithConfig(config Config) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}

	// New jpillora/ipfilter instance
	f, err := ipfilter.New(ipfilter.Options{
		AllowedIPs:     config.WhiteList,
		BlockedIPs:     config.BlackList,
		BlockByDefault: config.BlockByDefault,
		Logger:         nil,
	})
	if err != nil {
		panic(err)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			ip := c.RealIP()
			if ip == "" {
				ip, _, err = net.SplitHostPort(c.Request().RemoteAddr)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err.Error())
				}
			}
			if !f.Allowed(ip) {
				return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("IP address %s not allowed", ip))
			}
			return next(c)
		}
	}
}