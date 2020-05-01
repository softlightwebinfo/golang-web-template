package functions

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func writeCookie(c echo.Context, name string, value string)  {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
}

func readCookie(c echo.Context, name string ) (*http.Cookie,error) {
	cookie, err := c.Cookie(name)
	if err != nil {
		return nil, err
	}
	return cookie,nil
}