package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"trusteeTestTask/handlers"
)

func globalErrHandler(err error, ctx echo.Context) {
	httpErr := err.(*echo.HTTPError)
	var str string
	if e, ok := httpErr.Message.(error); ok {
		str = e.Error()
	} else {
		str = httpErr.Message.(string)
	}
	ctx.Logger().Errorj(log.JSON{
		"error": str,
		"ip":    ctx.RealIP(),
		"code":  httpErr.Code,
		"path":  ctx.Path(),
	})
	_ = ctx.JSON(httpErr.Code, echo.Map{
		"error": str,
	})
}

func main() {
	e := echo.New()
	e.HTTPErrorHandler = globalErrHandler

	h := handlers.New()

	e.GET("/authenticate", h.Authenticate)

	api := e.Group("/api")
	api.Use(h.TokenMiddleware)
	api.GET("/showtime", h.Showtime)
	api.PATCH("/showtime/:id", h.BookShowtime)
	api.DELETE("/showtime/:id", h.CancelBooking)

	//log.Fatal(e.Start(":" + os.Getenv("PORT")))
	log.Fatal(e.Start(":3000"))
}
