package ipfilter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	ipfilter "github.com/crazy-max/echo-ipfilter"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	a := assert.New(t)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1"
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	h := ipfilter.Middleware()(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	a.Error(h(c))
}

func TestMiddlewareWithConfig(t *testing.T) {
	a := assert.New(t)

	testTable := []struct {
		Config ipfilter.Config
		Ip     string
		Expect error
	}{
		{
			Config: ipfilter.Config{
				WhiteList: []string{
					"123.123.123.123",
				},
				BlockByDefault: true,
			},
			Ip: "223.123.123.123:1234",
			// Blocked by WhiteList.
			Expect: echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("IP address %s not allowed", "223.123.123.123")),
		},
		{
			Config: ipfilter.Config{
				WhiteList: []string{
					"223.123.123.123",
				},
			},
			Ip: "223.123.123.123:1234",
			// Not blocked.
			Expect: nil,
		},
		{
			Config: ipfilter.Config{
				WhiteList: []string{
					"223.123.123.0/24",
				},
			},
			Ip: "223.123.123.230:1234",
			// Not blocked.
			Expect: nil,
		},
		{
			Config: ipfilter.Config{
				WhiteList: []string{
					"223.123.123.0/24",
				},
				BlackList: []string{
					"223.123.123.230",
				},
			},
			Ip: "223.123.123.230:1234",
			// Blocked by BlackList.
			Expect: echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("IP address %s not allowed", "223.123.123.230")),
		},
		{
			// It will be DefaultIpFilterConfig.
			Config: ipfilter.Config{},
			// Allowed by default WhiteList.
			Ip:     "123.123.123.123:1234",
			Expect: nil,
		},
		{
			Config: ipfilter.Config{
				WhiteList: []string{
					"10.1.2.0/24",
					"10.1.4.0/24",
				},
				BlockByDefault: true,
			},
			Ip: "10.1.4.1:80",
			// Allowed by WhiteList.
			Expect: nil,
		},
	}

	e := echo.New()

	for idx, testCase := range testTable {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = testCase.Ip
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)
		h := ipfilter.MiddlewareWithConfig(testCase.Config)(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		})

		switch err := h(c); err.(type) {
		case nil:
			a.EqualValues(testCase.Expect, err, "testTable[%d]", idx)
		case *echo.HTTPError:
			a.EqualValues(testCase.Expect, err, "testTable[%d]", idx)
		}
	}
}
