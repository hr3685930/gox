package echo

import (
    "bytes"
    "fmt"
    "github.com/ddliu/go-httpclient"
    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo"
    "github.com/pkg/errors"
    "github.com/spf13/viper"
    "io/ioutil"
    "net/http"
)

type HttpError struct {
    HttpCode int
    Code     int    `json:"code"`
    Msg      string `json:"message"`
    Stack    []byte
}

func (h *HttpError) Error() string {
    return h.Msg
}

var dontReport = []int{
    http.StatusUnauthorized,
    http.StatusForbidden,
    http.StatusMethodNotAllowed,
    http.StatusUnsupportedMediaType,
    http.StatusUnprocessableEntity,
}

func CustomHTTPErrorHandler(err error, c echo.Context) {

    var (
        code = http.StatusInternalServerError
        msg  interface{}
    )
    var body string
    stack := fmt.Sprintf("%+v\n", errors.New("debug"))
    if he, ok := err.(*echo.HTTPError); ok {
        code = he.Code
        msg = echo.Map{
            "message": he.Message,
        }
        if he.Internal != nil {
            err = fmt.Errorf("%v, %v", err, he.Internal)
        }
    } else if he, ok := err.(*HttpError); ok {
        code = he.HttpCode
        stack = string(he.Stack)
        msg = echo.Map{
            "code":    he.Code,
            "message": he.Msg,
        }
        stack = string(he.Stack)
    } else if c.Echo().Debug {
        msg = err.Error()
    } else if errs, ok := err.(validator.ValidationErrors); ok {
        code = http.StatusUnprocessableEntity
        var details []string
        trans, _ := uni.GetTranslator("zh")
        for _, e := range errs {
            // can translate each errors one at a time.
            details = append(details, e.Translate(trans))
        }
        msg = echo.Map{
            "code":    4422,
            "message": "validation_failed",
            "detail":  details,
        }
    } else {
        msg = echo.Map{
            "message": err.Error(),
        }
    }

    isDontReport := false
    for _, value := range dontReport {
        if value == code {
            isDontReport = true
        }
    }

    errUrl := viper.GetString("error.report")
    if errUrl != "" && !isDontReport {
        bodyBytes, _ := ioutil.ReadAll(c.Request().Body)
        body = string(bodyBytes)

        request := map[string]interface{}{
            "url":     c.Request().Host + c.Request().RequestURI,
            "method":  c.Request().Method,
            "headers": c.Request().Header,
            "params":  body,
        }

        app := map[string]string{
            "name":        viper.GetString("app.name"),
            "environment": viper.GetString("app.env"),
        }

        exception := map[string]interface{}{
            "code":  code,
            "trace": stack,
        }

        option := map[string]interface{}{
            "error_type": "api_error",
            "app":        app,
            "exception":  exception,
            "request":    request,
        }
        go func() {
            _, _ = httpclient.Begin().PostJson(errUrl, option)
        }()
    }

    // Send response
    if !c.Response().Committed {
        if c.Request().Method == http.MethodHead { // Issue #608
            err = c.NoContent(code)
        } else {
            err = c.JSON(code, msg)
        }
        if err != nil {
            c.Echo().Logger.Error(err)
        }
    }
}

type CustomBinder struct{}

func NewCustomBinder() *CustomBinder {
    return &CustomBinder{}
}

func (cb *CustomBinder) Bind(i interface{}, c echo.Context) (err error) {
    bodyBytes, _ := ioutil.ReadAll(c.Request().Body)
    // Restore the io.ReadCloser to its original state
    c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
    // You may use default binder
    db := new(echo.DefaultBinder)
    err = db.Bind(i, c)
    c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
    if err != echo.ErrUnsupportedMediaType {
        return
    }
    // Define your custom implementation
    return
}
