package echo

import (
    "github.com/labstack/echo"
    "net/http"
)

type EchoHTTP struct {

}

func (*EchoHTTP) HTTP(s *http.Server,route func(echo2 *echo.Echo)) error {
    e := echo.New()
    echo.NotFoundHandler = func(c echo.Context) error {
        return echo.ErrMethodNotAllowed
    }
    e.HTTPErrorHandler = CustomHTTPErrorHandler
    e.Validator = NewCustomValidator()
    e.Binder = NewCustomBinder()

    route(e)
    return e.StartServer(s)
}
