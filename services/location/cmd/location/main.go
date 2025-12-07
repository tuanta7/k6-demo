package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tuanta7/k6-demo/services/location/pkg/zapx"
)

func main() {
	logger, err := zapx.NewZapLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
