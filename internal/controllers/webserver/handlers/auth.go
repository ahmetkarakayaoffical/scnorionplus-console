package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/views/register_views"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type MyCustomClaims struct {
	jwt.RegisteredClaims
}

func (h *Handler) Auth(c echo.Context) error {
	if h.ReverseProxyAuthPort != "" {
		u, err := url.Parse(c.Request().Referer())
		if err != nil {
			return views.ErrorPage(strconv.Itoa(http.StatusBadRequest), "could not parse url").Render(c.Request().Context(), c.Response().Writer)
		}
		return c.Redirect(http.StatusFound, fmt.Sprintf("https://%s:%s/auth", u.Hostname(), h.ReverseProxyAuthPort))
	} else {
		return c.Redirect(http.StatusFound, fmt.Sprintf("https://%s:%s/auth", h.ServerName, h.AuthPort))
	}

}

func (h *Handler) ConfirmEmail(c echo.Context) error {
	tokenString := c.Param("token")

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.JWTKey), nil
	})

	if err != nil {
		return views.ErrorPage(strconv.Itoa(http.StatusBadRequest), err.Error()).Render(c.Request().Context(), c.Response().Writer)
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok {
		if time.Now().After(claims.ExpiresAt.Time) {
			return views.ErrorPage(strconv.Itoa(http.StatusBadRequest), "token has expired, please contact your administrator to request a new confirmation email").Render(c.Request().Context(), c.Response().Writer)
		}

		user, err := h.Model.GetUserById(claims.ID)
		if err != nil {
			return views.ErrorPage(strconv.Itoa(http.StatusBadRequest), err.Error()).Render(c.Request().Context(), c.Response().Writer)
		}

		if user.EmailVerified {
			return views.ErrorPage(strconv.Itoa(http.StatusBadRequest), "you've already confirmed your email").Render(c.Request().Context(), c.Response().Writer)
		}

		if err := h.Model.ConfirmEmail(user.ID); err != nil {
			return views.ErrorPage(strconv.Itoa(http.StatusInternalServerError), err.Error()).Render(c.Request().Context(), c.Response().Writer)
		}

		return RenderView(c, register_views.RegisterIndex(register_views.EmailConfirmed()))

	} else {
		return views.ErrorPage(strconv.Itoa(http.StatusBadRequest), "unknown claims type, cannot proceed").Render(c.Request().Context(), c.Response().Writer)
	}
}
