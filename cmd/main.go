package main

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"trusteeTestTask/handlers"
)

func globalErrHandler(err error, ctx echo.Context) {
	code := 500
	msg := err.Error()
	if e := new(echo.HTTPError); errors.As(err, &e) {
		switch e.Message.(type) {
		case error:
			msg = e.Message.(error).Error()
		case string:
			msg = e.Message.(string)
		}
		code = e.Code
	}
	ctx.Logger().Errorj(log.JSON{
		"error": msg,
		"ip":    ctx.RealIP(),
		"code":  code,
		"path":  ctx.Request().RequestURI,
	})
	_ = ctx.JSON(code, echo.Map{
		"error": msg,
	})
}

func main() {
	e := echo.New()
	e.HTTPErrorHandler = globalErrHandler

	h := handlers.New()

	e.GET("/authenticate", h.Authenticate)

	api := e.Group("/api")
	api.Use(h.TokenMiddleware)
	api.GET("/showtime", h.Showtimes)
	api.GET("/showtime/:id", h.Showtime)
	api.PATCH("/showtime/:id", h.BookShowtime)
	api.DELETE("/showtime/:id", h.CancelBooking)

	// we can also do go:embed if our src folder is not a blob storage volume
	//e.Static("/src/", os.Getenv("IMAGE_VOLUME_PATH"))
	e.Static("/src/", "D:/Pics")

	//log.Fatal(e.Start(":" + os.Getenv("PORT")))
	log.Fatal(e.Start(":3000"))
}
