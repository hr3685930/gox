package gin

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
	"runtime/debug"
)

var (
	InternalError = NewError(http.StatusInternalServerError, 5500, http.StatusText(http.StatusInternalServerError), debug.Stack())
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

func NewError(statusCode, code int, msg string, stack []byte) *HttpError {
	return &HttpError{
		HttpCode: statusCode,
		Code:     code,
		Msg:      msg,
		Stack:    stack,
	}
}

type HTTPErrorReport func(HTTPCode int, response gin.H, stack string, c *gin.Context)

func ErrHandler(errorReport HTTPErrorReport) gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Set("jsonBody", string(bodyBytes))
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
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
					stack = string(debug.Stack())
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
