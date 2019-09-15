<p align="center"><img width="156" src=".res/echo-ipfilter.png"></p>

<p align="center">
  <a href="https://github.com/crazy-max/echo-ipfilter/actions"><img src="https://github.com/crazy-max/echo-ipfilter/workflows/build/badge.svg" alt="Build Status"></a>
  <a href="https://goreportcard.com/report/github.com/crazy-max/echo-ipfilter"><img src="https://goreportcard.com/badge/github.com/crazy-max/echo-ipfilter?style=flat-square" alt="Go Report"></a>
  <a href="https://www.codacy.com/app/crazy-max/echo-ipfilter"><img src="https://img.shields.io/codacy/grade/99e6d0f21bd6475187823203da6fce63/master.svg?style=flat-square" alt="Code Quality"></a>
  <a href="https://codecov.io/gh/crazy-max/echo-ipfilter"><img src="https://img.shields.io/codecov/c/github/crazy-max/echo-ipfilter/master.svg?style=flat-square" alt="Codecov"></a>
  <a href="https://www.patreon.com/crazymax"><img src="https://img.shields.io/badge/donate-patreon-f96854.svg?logo=patreon&style=flat-square" alt="Support me on Patreon"></a>
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

## How can I help ?

All kinds of contributions are welcome :raised_hands:!<br />
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:<br />
But we're not gonna lie to each other, I'd rather you buy me a beer or two :beers:!

[![Support me on Patreon](.res/patreon.png)](https://www.patreon.com/crazymax) 
[![Paypal Donate](.res/paypal.png)](https://www.paypal.me/crazyws)

## License

MIT. See `LICENSE` for more details.
