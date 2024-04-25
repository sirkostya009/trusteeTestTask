package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"time"
)

func (h *Handlers) TokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth, err := c.Cookie("google-session")
		if err != nil {
			return echo.NewHTTPError(401, "no auth")
		}
		token, err := jwt.Parse(auth.Value, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return h.secret, nil
		})
		if err != nil {
			return echo.NewHTTPError(401, err)
		}

		claims := token.Claims.(jwt.MapClaims)
		exp, err := claims.GetExpirationTime()
		if err != nil {
			return echo.NewHTTPError(401, err)
		}
		if exp.Time.Before(time.Now()) {
			return echo.NewHTTPError(401, "token expired")
		}

		ctx := c.Request().Context()
		client := h.config.Client(ctx, &oauth2.Token{
			RefreshToken: claims["ref"].(string),
			AccessToken:  claims["acs"].(string),
			TokenType:    claims["type"].(string),
			Expiry:       exp.Time,
		})
		srvc, err := calendar.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			return echo.NewHTTPError(500, err)
		}
		c.Set("service", srvc)

		return next(c)
	}
}
