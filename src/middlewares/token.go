package middlewares

import (
	"github.com/labstack/echo/v4"
)

//IsAdmin middleware evalutes if the user is admin - super admin
func Tokenizing(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		headertoken := c.Request().Header.Get("Authorization")
		c.Set("Authorization", headertoken)
		return next(c)
	}
}
