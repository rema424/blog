package sandbox

import "github.com/labstack/echo/v4"

func HelloHandler(c echo.Context) error {
	// return c.String(200, "Hello, World!")
	return c.Render(200, "top.html", nil)
}
