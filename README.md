<p align="center"><img width="156" src="https://raw.githubusercontent.com/crazy-max/echo-ipfilter/master/.github/echo-ipfilter.png"></p>

<p align="center">
  <a href="https://github.com/crazy-max/echo-ipfilter/actions?workflow=test"><img src="https://img.shields.io/github/workflow/status/crazy-max/echo-ipfilter/test?label=test&logo=github&style=flat-square" alt="Test workflow"></a>
  <a href="https://goreportcard.com/report/github.com/crazy-max/echo-ipfilter"><img src="https://goreportcard.com/badge/github.com/crazy-max/echo-ipfilter?style=flat-square" alt="Go Report"></a>
  <a href="https://www.codacy.com/app/crazy-max/echo-ipfilter"><img src="https://img.shields.io/codacy/grade/99e6d0f21bd6475187823203da6fce63/master.svg?style=flat-square" alt="Code Quality"></a>
  <a href="https://codecov.io/gh/crazy-max/echo-ipfilter"><img src="https://img.shields.io/codecov/c/github/crazy-max/echo-ipfilter?logo=codecov&style=flat-square" alt="Codecov"></a>
  <br /><a href="https://github.com/sponsors/crazy-max"><img src="https://img.shields.io/badge/sponsor-crazy--max-181717.svg?logo=github&style=flat-square" alt="Become a sponsor"></a>
  <a href="https://www.paypal.me/crazyws"><img src="https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal&style=flat-square" alt="Donate Paypal"></a>
</p>

## About

Middleware that provides ipfilter support for [echo](https://echo.labstack.com) framework backed by [jpillora/ipfilter](https://github.com/jpillora/ipfilter).

## Installation

```
go get github.com/crazy-max/echo-ipfilter
```

## Example

```go
package main

import (
	"net/http"

	ipfilter "github.com/crazy-max/echo-ipfilter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(ipfilter.MiddlewareWithConfig(ipfilter.Config{
		Skipper: middleware.DefaultSkipper,
		WhiteList: []string{
			"10.1.1.0/24",
			"10.1.2.0/24",
		},
		BlockByDefault: true,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
```

## How can I help?

All kinds of contributions are welcome :raised_hands:! The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon: You can also support this project by [**becoming a sponsor on GitHub**](https://github.com/sponsors/crazy-max) :clap: or by making a [Paypal donation](https://www.paypal.me/crazyws) to ensure this journey continues indefinitely! :rocket:

Thanks again for your support, it is much appreciated! :pray:

## License

MIT. See `LICENSE` for more details.
