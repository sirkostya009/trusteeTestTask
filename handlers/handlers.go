package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"net/http"
	"os"
	"strconv"
	"time"
	"trusteeTestTask/service"
)

type Handlers struct {
	config  *oauth2.Config
	secret  []byte
	service *service.CinemaService
}

func New() *Handlers {
	h := &Handlers{}

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("No creds %v", err)
	}

	h.config, err = google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("oopsie woopsie")
	}

	//h.secret = []byte(os.Getenv("secret"))
	h.secret = []byte("secret")

	h.service = service.NewCinemaService()

	return h
}

func (h *Handlers) Authenticate(c echo.Context) error {
	authCode := c.QueryParam("code")
	if authCode == "" {
		return c.Redirect(302, h.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline))
	}

	tok, err := h.config.Exchange(c.Request().Context(), authCode)
	if err != nil {
		return echo.NewHTTPError(400, err)
	}

	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  jwt.NumericDate{Time: tok.Expiry},
		"ref":  tok.RefreshToken,
		"acs":  tok.AccessToken,
		"type": tok.TokenType,
	}).SignedString(h.secret)
	if err != nil {
		return echo.NewHTTPError(500, err)
	}

	c.SetCookie(&http.Cookie{
		Name:     "google-session",
		Value:    s,
		Path:     "/api/",
		Expires:  tok.Expiry,
		HttpOnly: true,
		Secure:   true,
	})

	return c.String(200, "goggle-session cookie is set!")
}

func (h *Handlers) getEvents(srvc *calendar.Service) (*calendar.Events, error) {
	return srvc.Events.List("primary").
		TimeMin(time.Now().Format(time.RFC3339)).
		TimeMax(h.service.MaxShowtimeEndTime().Format(time.RFC3339)).
		Do()
}

func (h *Handlers) Showtimes(c echo.Context) error {
	events, err := h.getEvents(c.Get("service").(*calendar.Service))
	if err != nil {
		return echo.NewHTTPError(500, err)
	}

	showtimeList := h.service.FilterShowtime(events.Items)

	return c.JSON(200, showtimeList)
}

func (h *Handlers) Showtime(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(400, err)
	}

	showtime, err := h.service.GetShowtime(id)
	if err != nil {
		return echo.NewHTTPError(404, "showtime not found")
	}

	return c.JSON(200, showtime.ToDto(c.RealIP()))
}

func (h *Handlers) BookShowtime(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(400, err)
	}

	if h.service.IsBooked(id, c.RealIP()) {
		return echo.NewHTTPError(400, "already booked")
	}

	showtime, err := h.service.GetShowtime(id)
	if err != nil {
		return echo.NewHTTPError(404, "showtime not found")
	}

	events, err := h.getEvents(c.Get("service").(*calendar.Service))
	if err != nil {
		return echo.NewHTTPError(500, err)
	}

	if h.service.IsShowtimeOverlapping(events.Items, showtime) {
		return echo.NewHTTPError(400, "showtime is overlapping")
	}

	srvc := c.Get("service").(*calendar.Service)

	event, err := srvc.Events.Insert("primary", &calendar.Event{
		Status:      "tentative",
		Summary:     showtime.Title,
		Location:    "Адреса кіношніка",
		Description: "Похід в кіно!",
		//Creator: &calendar.EventCreator{
		//	DisplayName: "Кіношнік",
		//},
		Start: &calendar.EventDateTime{
			DateTime: showtime.BeginTime.Format(time.RFC3339),
			TimeZone: "Europe/Kiev",
		},
		End: &calendar.EventDateTime{
			DateTime: showtime.EndTime.Format(time.RFC3339),
			TimeZone: "Europe/Kiev",
		},
		Locked: true,
	}).Do()
	if err != nil {
		return echo.NewHTTPError(500, err)
	}

	_ = h.service.BookShowtime(id, c.RealIP(), event.Id)

	return c.String(200, event.Id)
}

func (h *Handlers) CancelBooking(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(500, err)
	}

	if !h.service.IsBooked(id, c.RealIP()) {
		return echo.NewHTTPError(400, "nothing to cancel")
	}

	eventId := h.service.CancelBooking(id, c.RealIP())

	if eventId != "" {
		srvc := c.Get("service").(*calendar.Service)
		err = srvc.Events.Delete("primary", eventId).Do()
		if err != nil {
			c.Logger().Error(err)
		}
	}

	return c.JSON(200, map[string]any{
		"eventDeleted": eventId != "" && err == nil,
	})
}
