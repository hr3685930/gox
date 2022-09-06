package gin

import (
	"bytes"
	"fmt"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	sentinelPlugin "github.com/sentinel-group/sentinel-go-adapters/gin"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	InternalError = NewError(http.StatusInternalServerError, 5500, http.StatusText(http.StatusInternalServerError))
)

type HttpError struct {
	HttpCode int    `json:"-"`
	Code     int    `json:"code"`
	Msg      string `json:"message"`
	Stack    []byte `json:"-"`
}

func (h *HttpError) Error() string {
	return h.Msg
}

func (h *HttpError) GetStack() string {
	return string(h.Stack)
}

func NewError(statusCode, code int, msg string) *HttpError {
	return &HttpError{
		HttpCode: statusCode,
		Code:     code,
		Msg:      msg,
		Stack:    []byte(fmt.Sprintf("%+v\n", errors.New(msg))),
	}
}

type HTTPErrorReport func(HTTPCode int, response gin.H, stack string, c *gin.Context)

func ErrHandler(errorReport HTTPErrorReport) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil  {
			bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
			c.Set("jsonBody", string(bodyBytes))
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		c.Next()
		if length := len(c.Errors); length > 0 {
			err := c.Errors[length-1].Err
			response := gin.H{}
			if err != nil {
				var HTTPCode = http.StatusInternalServerError
				var stack string
				if e, ok := err.(*HttpError); ok {
					HTTPCode = e.HttpCode
					response["code"] = e.Code
					response["message"] = e.Msg
					stack = string(e.Stack)
				} else if e, ok := err.(validator.ValidationErrors); ok {
					HTTPCode = http.StatusUnprocessableEntity
					response["code"] = 4422
					response["message"] = "validation_failed"
					response["detail"] = Translate(e)
					stack = fmt.Sprintf("%+v\n", errors.New("validation_failed"))
				} else {
					response["code"] = InternalError.Code
					response["message"] = InternalError.Msg
					stack = string(InternalError.Stack)
				}

				// error report
				errorReport(HTTPCode, response, stack, c)

				c.JSON(HTTPCode, response)
				return
			}
		}

	}
}

// TimeoutMiddleware timeout middleware wraps the request context with a timeout
func TimeoutMiddleware(duration time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(duration),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"code": 5504, "message": "status gateway timeout"})
		}),
	)
}

// GovernanceMiddleware service governance middleware
func GovernanceMiddleware(fn func(ctx *gin.Context)) gin.HandlerFunc {
	return sentinelPlugin.SentinelMiddleware(
		sentinelPlugin.WithBlockFallback(fn),
	)
}
