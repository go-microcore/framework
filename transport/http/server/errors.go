package server // import "go.microcore.dev/framework/transport/http/server"

import (
	"github.com/valyala/fasthttp"
	
	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/errors"
)

var ErrorStatusCodeMap = map[error]int{
	// Client
	errors.ErrBadRequest:                   fasthttp.StatusBadRequest,
	errors.ErrUnauthorized:                 fasthttp.StatusUnauthorized,
	errors.ErrPaymentRequired:              fasthttp.StatusPaymentRequired,
	errors.ErrForbidden:                    fasthttp.StatusForbidden,
	errors.ErrNotFound:                     fasthttp.StatusNotFound,
	errors.ErrMethodNotAllowed:             fasthttp.StatusMethodNotAllowed,
	errors.ErrConflict:                     fasthttp.StatusConflict,
	errors.ErrGone:                         fasthttp.StatusGone,
	errors.ErrRequestEntityTooLarge:        fasthttp.StatusRequestEntityTooLarge,
	errors.ErrUnsupportedMediaType:         fasthttp.StatusUnsupportedMediaType,
	errors.ErrRequestedRangeNotSatisfiable: fasthttp.StatusRequestedRangeNotSatisfiable,
	errors.ErrTeapot:                       fasthttp.StatusTeapot,
	errors.ErrUnprocessableEntity:          fasthttp.StatusUnprocessableEntity,
	errors.ErrLocked:                       fasthttp.StatusLocked,
	errors.ErrFailedDependency:             fasthttp.StatusFailedDependency,
	errors.ErrPreconditionRequired:         fasthttp.StatusPreconditionRequired,
	errors.ErrTooManyRequests:              fasthttp.StatusTooManyRequests,
	errors.ErrUnavailableForLegalReasons:   fasthttp.StatusUnavailableForLegalReasons,
	// Server
	errors.ErrInternalServerError: fasthttp.StatusInternalServerError,
	errors.ErrNotImplemented:      fasthttp.StatusNotImplemented,
	errors.ErrBadGateway:          fasthttp.StatusBadGateway,
	errors.ErrServiceUnavailable:  fasthttp.StatusServiceUnavailable,
	errors.ErrGatewayTimeout:      fasthttp.StatusGatewayTimeout,
	errors.ErrLoopDetected:        fasthttp.StatusLoopDetected,
}
