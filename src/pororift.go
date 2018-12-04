package main

import (
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Host struct {
		Echo *echo.Echo
	}
)

func main() {
	// Host
	hosts := map[string]*Host{}

	// API Proxy
	api := echo.New()
	api.Use(middleware.Logger())
	api.Use(middleware.Recover())

	hosts["api.localhost:3000"] = &Host{api}

	url1, err := url.Parse("http://backend:3001")
	if err != nil {
		api.Logger.Fatal(err)
	}

	targets := []*middleware.ProxyTarget{
		{
			URL: url1,
		},
	}

	api.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))

	// Website
	site := echo.New()
	site.Use(middleware.Logger())
	site.Use(middleware.Recover())

	hosts["localhost:3000"] = &Host{site}

	site.Static("/", "./dist/pororift-client/")
	site.File("/", "./dist/pororift-client/index.html")
	// site.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Website")
	// })

	// Server
	e := echo.New()

	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]

		if host == nil {
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return
	})

	e.Logger.Fatal(e.Start(":3000"))
}
