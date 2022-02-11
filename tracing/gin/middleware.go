package gin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"io"
	"net/http"
)

func OpenTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		var parentSpan opentracing.Span
		if !opentracing.IsGlobalTracerRegistered() {
			c.Next()
			return
		}

		spCtx, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header),
		)
		if err != nil {
			parentSpan = opentracing.GlobalTracer().StartSpan(
				c.Request.URL.Path,
				ext.SpanKindRPCServer,
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			)
		} else {
			parentSpan = opentracing.GlobalTracer().StartSpan(
				c.Request.URL.Path,
				ext.RPCServerOption(spCtx),
				ext.SpanKindRPCServer,
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			)
		}
		_ = opentracing.GlobalTracer().Inject(
			parentSpan.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header),
		)
		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), parentSpan))

		c.Next()

		ext.HTTPMethod.Set(parentSpan, c.Request.Method)
		ext.HTTPUrl.Set(parentSpan, c.Request.URL.String())
		ext.HTTPStatusCode.Set(parentSpan, uint16(c.Writer.Status()))

		if length := len(c.Errors); length > 0 {
			err := c.Errors[length-1].Err
			finishSpan(parentSpan, err)
		} else if length == 0 && c.Writer.Status() >= http.StatusInternalServerError {
			finishSpan(parentSpan, errors.New("返回状态码不正确"))
		} else {
			parentSpan.Finish()
		}
	}
}

func finishSpan(span opentracing.Span, err error) {
	if err != nil && err != io.EOF {
		ext.Error.Set(span, true)
		span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
	}
	span.Finish()
}
