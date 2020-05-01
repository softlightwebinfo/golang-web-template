package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func HomeRoute(c echo.Context) error {
	return c.Render(http.StatusOK, "template.html", map[string]interface{}{
		"name": "Dolly!",
	})
}