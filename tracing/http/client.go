package http

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/parnurzeal/gorequest"
	"io"
	"net/http"
)

type httpClient struct {
	*gorequest.SuperAgent
}

func NewHttpClient() *httpClient {
	return &httpClient{gorequest.New()}
}

func (h *httpClient) TraceEnd(ctx context.Context) (res *http.Response, body string, errs []error) {
	if !opentracing.IsGlobalTracerRegistered() {
		return h.End()
	}
	respSpan, _ := opentracing.StartSpanFromContext(
		ctx,
		"http_client",
		ext.SpanKindRPCClient,
		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
		opentracing.Tag{Key: string(ext.HTTPUrl), Value: h.Url},
		opentracing.Tag{Key: string(ext.HTTPMethod), Value: h.Method},
	)

	_ = opentracing.GlobalTracer().Inject(
		respSpan.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(h.Header),
	)

	res, body, errs = h.End()

	if length := len(errs); length > 0 {
		err := errs[length-1]
		finishSpan(respSpan, err)
	} else if length == 0 && res.StatusCode >= http.StatusMultipleChoices {
		ext.HTTPStatusCode.Set(respSpan, uint16(res.StatusCode))
		finishSpan(respSpan, errors.New("返回状态码不正确"))
	} else {
		ext.HTTPStatusCode.Set(respSpan, uint16(res.StatusCode))
		respSpan.Finish()
	}

	return res, body, errs
}


func (h *httpClient) TraceEndByte(ctx context.Context) (res *http.Response, body []byte, errs []error) {
	if !opentracing.IsGlobalTracerRegistered() {
		return h.EndBytes()
	}
	respSpan, _ := opentracing.StartSpanFromContext(
		ctx,
		"http_client",
		ext.SpanKindRPCClient,
		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
		opentracing.Tag{Key: string(ext.HTTPUrl), Value: h.Url},
		opentracing.Tag{Key: string(ext.HTTPMethod), Value: h.Method},
	)

	_ = opentracing.GlobalTracer().Inject(
		respSpan.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(h.Header),
	)

	res, body, errs = h.EndBytes()

	if length := len(errs); length > 0 {
		err := errs[length-1]
		finishSpan(respSpan, err)
	} else if length == 0 && res.StatusCode >= http.StatusMultipleChoices {
		ext.HTTPStatusCode.Set(respSpan, uint16(res.StatusCode))
		finishSpan(respSpan, errors.New("返回状态码不正确"))
	} else {
		ext.HTTPStatusCode.Set(respSpan, uint16(res.StatusCode))
		respSpan.Finish()
	}

	return res, body, errs
}


func finishSpan(span opentracing.Span, err error) {
	if err != nil && err != io.EOF {
		ext.Error.Set(span, true)
		span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
	}
	span.Finish()
}
