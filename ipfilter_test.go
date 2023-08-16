package ipfilter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	ipfilter "github.com/crazy-max/echo-ipfilter"
	jpillorafilter "github.com/jpillora/ipfilter"
	"github.com/labstack/echo/v4"
)

func TestMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1"
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	h := ipfilter.Middleware()(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	if err := h(c); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestMiddlewareWithConfig(t *testing.T) {
	cases := []struct {
		name   string
		config ipfilter.Config
		ip     string
		err    error
	}{
		{
			name: "blocked with whitelist",
			config: ipfilter.Config{
				WhiteList: []string{
					"123.123.123.123",
				},
				BlockByDefault: true,
			},
			ip:  "223.123.123.123:1234",
			err: echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("IP address %s not allowed", "223.123.123.123")),
		},
		{
			name: "not blocked",
			config: ipfilter.Config{
				WhiteList: []string{
					"223.123.123.123",
				},
			},
			ip: "223.123.123.123:1234",
			// Not blocked.
			err: nil,
		},
		{
			name: "not blocked cidr",
			config: ipfilter.Config{
				WhiteList: []string{
					"223.123.123.0/24",
				},
			},
			ip:  "223.123.123.230:1234",
			err: nil,
		},
		{
			name: "blocked by blacklist",
			config: ipfilter.Config{
				WhiteList: []string{
					"223.123.123.0/24",
				},
				BlackList: []string{
					"223.123.123.230",
				},
			},
			ip:  "223.123.123.230:1234",
			err: echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("IP address %s not allowed", "223.123.123.230")),
		},
		{
			name:   "allowed by default whitelist",
			config: ipfilter.Config{},
			ip:     "123.123.123.123:1234",
			err:    nil,
		},
		{
			name: "allowed by whitelist",
			config: ipfilter.Config{
				WhiteList: []string{
					"10.1.2.0/24",
					"10.1.4.0/24",
				},
				BlockByDefault: true,
			},
			ip:  "10.1.4.1:80",
			err: nil,
		},
		{
			name: "dynamically allowed by whitelist",
			config: ipfilter.Config{
				WhiteList: []string{
					"10.1.2.0/24",
					// this will be dynamically added "10.1.4.0/24",
				},
				BlockByDefault: true,
				CreatedFilter: func(filter *jpillorafilter.IPFilter) {
					filter.AllowIP("10.1.4.0/24")
				},
			},
			ip:  "10.1.4.1:80",
			err: nil,
		},
		{
			name: "dynamically allowed by whitelist",
			config: ipfilter.Config{
				WhiteList: []string{
					"10.1.2.0/24", // will be dynamicaly blocked
					"10.1.4.0/24",
				},
				BlockByDefault: true,
				CreatedFilter: func(filter *jpillorafilter.IPFilter) {
					filter.BlockIP("10.1.2.0/24")
				},
			},
			ip:  "10.1.2.7:80",
			err: echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("IP address %s not allowed", "10.1.2.7")),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = tt.ip
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			h := ipfilter.MiddlewareWithConfig(tt.config)(func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			})

			switch err := h(c); err.(type) {
			case nil:
				if tt.err != nil {
					t.Errorf("expected error %v, got nil", tt.err)
				}
			case *echo.HTTPError:
				if tt.err.Error() != err.Error() {
					t.Errorf("expected error %v, got %v", tt.err, err)
				}
			}
		})
	}
}
