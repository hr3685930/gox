package echo

import (
	"github.com/labstack/echo"
)

type EchoHTTP struct {
	E     *echo.Echo
	Debug bool
}

func NewEchoHTTP(debug bool) *EchoHTTP {
	return &EchoHTTP{Debug: debug}
}

func (ec *EchoHTTP) LoadRoute(route func(echo2 *echo.Echo)) error {
	e := echo.New()
	echo.NotFoundHandler = func(c echo.Context) error {
		return echo.ErrMethodNotAllowed
	}
	e.HideBanner = true
	//e.HTTPErrorHandler = CustomHTTPErrorHandler
	e.Validator = NewCustomValidator()
	e.Binder = NewCustomBinder()
	route(e)
	ec.E = e
	return nil
}

func (e *EchoHTTP) Run(addr string) error {
	return e.E.Start(addr)
}
