package server // import "go.microcore.dev/framework/transport/http/server"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/trace"

	_ "go.microcore.dev/framework"
)

type RequestContext struct {
	*fasthttp.RequestCtx
}

func (c *RequestContext) Write(p []byte) (int, error) {
	c.SetTraceIdHeader()
	return c.RequestCtx.Write(p)
}

func (c *RequestContext) WriteString(s string) (int, error) {
	c.SetTraceIdHeader()
	return c.RequestCtx.WriteString(s)
}

func (c *RequestContext) WriteError(err error) {
	for e, s := range ErrorStatusCodeMap {
		if errors.Is(err, e) {
			c.Error(err.Error(), s)
			c.SetTraceIdHeader()
			return
		}
	}
	c.Error(
		defaultResponseError.Error(),
		ErrorStatusCodeMap[defaultResponseError],
	)
	c.SetTraceIdHeader()
}

func (c *RequestContext) WriteJson(data any) error {
	c.SetTraceIdHeader()
	c.SetContentType("application/json; charset=utf-8")
	return json.NewEncoder(c).Encode(data)
}

func (c *RequestContext) WriteJsonWithStatusCode(statusCode int, data any) error {
	c.SetStatusCode(statusCode)
	if data == nil {
		c.SetTraceIdHeader()
		return nil
	}
	return c.WriteJson(data)
}

func (c *RequestContext) WriteStatusCode(code int) {
	c.SetTraceIdHeader()
	c.SetStatusCode(code)
}

func (c *RequestContext) ReadJsonBody(data any) error {
	return json.Unmarshal(c.Request.Body(), data)
}

func (c *RequestContext) UserValueBool(key any) (bool, error) {
	v := c.UserValue(key)
	if v == nil {
		return false, fmt.Errorf("key %v not found", key)
	}

	switch val := v.(type) {
	case bool:
		return val, nil
	case string:
		return strconv.ParseBool(val)
	case []byte:
		return strconv.ParseBool(string(val))
	case int:
		return val != 0, nil
	case int8:
		return val != 0, nil
	case int16:
		return val != 0, nil
	case int32:
		return val != 0, nil
	case int64:
		return val != 0, nil
	case uint:
		return val != 0, nil
	case uint8:
		return val != 0, nil
	case uint16:
		return val != 0, nil
	case uint32:
		return val != 0, nil
	case uint64:
		return val != 0, nil
	case float32:
		return val != 0, nil
	case float64:
		return val != 0, nil
	default:
		return false, fmt.Errorf("unsupported type %T for key %v", v, key)
	}
}

func (c *RequestContext) UserValueStr(key any) (string, error) {
	v := c.UserValue(key)
	if v == nil {
		return "", fmt.Errorf("key %v not found", key)
	}

	switch val := v.(type) {
	case string:
		return val, nil
	case []byte:
		return string(val), nil
	case fmt.Stringer:
		return val.String(), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func (c *RequestContext) UserValueInt(key any) (int, error) {
	v := c.UserValue(key)
	if v == nil {
		return 0, fmt.Errorf("key %v not found", key)
	}

	switch val := v.(type) {
	case int:
		return val, nil
	case int8:
		return int(val), nil
	case int16:
		return int(val), nil
	case int32:
		return int(val), nil
	case int64:
		return int(val), nil
	case uint:
		return int(val), nil
	case uint8:
		return int(val), nil
	case uint16:
		return int(val), nil
	case uint32:
		return int(val), nil
	case uint64:
		return int(val), nil
	case float32:
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		n, err := strconv.Atoi(val)
		if err != nil {
			return 0, fmt.Errorf("failed to parse string %q as int: %w", val, err)
		}
		return n, nil
	case []byte:
		n, err := strconv.Atoi(string(val))
		if err != nil {
			return 0, fmt.Errorf("failed to parse []byte %q as int: %w", val, err)
		}
		return n, nil
	case fmt.Stringer:
		n, err := strconv.Atoi(val.String())
		if err != nil {
			return 0, fmt.Errorf("failed to parse Stringer %q as int: %w", val.String(), err)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("unsupported type %T for key %v", v, key)
	}
}

func (c *RequestContext) UserValueUint(key any) (uint, error) {
	v := c.UserValue(key)
	if v == nil {
		return 0, fmt.Errorf("key %v not found", key)
	}

	switch val := v.(type) {
	case uint:
		return val, nil
	case uint8:
		return uint(val), nil
	case uint16:
		return uint(val), nil
	case uint32:
		return uint(val), nil
	case uint64:
		return uint(val), nil
	case int:
		return uint(val), nil
	case int8:
		return uint(val), nil
	case int16:
		return uint(val), nil
	case int32:
		return uint(val), nil
	case int64:
		return uint(val), nil
	case float32:
		return uint(val), nil
	case float64:
		return uint(val), nil
	case string:
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse string %q as uint: %w", val, err)
		}
		return uint(n), nil
	case []byte:
		n, err := strconv.ParseUint(string(val), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse []byte %q as uint: %w", val, err)
		}
		return uint(n), nil
	case fmt.Stringer:
		n, err := strconv.ParseUint(val.String(), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse Stringer %q as uint: %w", val.String(), err)
		}
		return uint(n), nil
	default:
		return 0, fmt.Errorf("unsupported type %T for key %v", v, key)
	}
}

func (c *RequestContext) SetTraceIdHeader() {
	span := trace.SpanContextFromContext(extractRequestContext(c.RequestCtx))
	if span.HasTraceID() {
		c.Response.Header.Set("X-Trace-Id", span.TraceID().String())
	}
}

func (c *RequestContext) GetHeaderStr(key string) string {
	return string(c.Request.Header.Peek(key))
}

func (c *RequestContext) GetIpAddr() string {
	if ip := c.Request.Header.Peek("X-Real-IP"); len(ip) > 0 {
		return string(ip)
	}
	if ip := c.Request.Header.Peek("X-Forwarded-For"); len(ip) > 0 {
		return string(ip)
	}
	if addr := c.Response.RemoteAddr(); addr != nil {
		ip := addr.String()
		if host, _, err := net.SplitHostPort(ip); err == nil {
			return host
		}
		return ip
	}
	return ""
}

func (c *RequestContext) GetBearerToken() (string, error) {
	authHeader := c.GetHeaderStr("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("missing bearer prefix")
	}
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
	return token, nil
}

func (c *RequestContext) GetContext() context.Context {
	return extractRequestContext(c.RequestCtx)
}

type fasthttpRequestCtxHeaderCarrier struct {
	ctx *fasthttp.RequestCtx
}

func (c fasthttpRequestCtxHeaderCarrier) Get(key string) string {
	return string(c.ctx.Request.Header.Peek(key))
}

func (c fasthttpRequestCtxHeaderCarrier) Set(key, value string) {
	c.ctx.Request.Header.Set(key, value)
}

func (c fasthttpRequestCtxHeaderCarrier) Keys() []string {
	keys := []string{}
	c.ctx.Request.Header.VisitAll(func(k, v []byte) {
		keys = append(keys, string(k))
	})
	return keys
}
