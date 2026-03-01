package http // import "go.microcore.dev/framework/transport/http"

import (
	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/transport"
)

var ErrStatusCodeMap = map[error]StatusCode{
	// Client
	transport.ErrBadRequest:                   StatusBadRequest,
	transport.ErrUnauthorized:                 StatusUnauthorized,
	transport.ErrPaymentRequired:              StatusPaymentRequired,
	transport.ErrForbidden:                    StatusForbidden,
	transport.ErrNotFound:                     StatusNotFound,
	transport.ErrMethodNotAllowed:             StatusMethodNotAllowed,
	transport.ErrConflict:                     StatusConflict,
	transport.ErrGone:                         StatusGone,
	transport.ErrRequestEntityTooLarge:        StatusPayloadTooLarge,
	transport.ErrUnsupportedMediaType:         StatusUnsupportedMediaType,
	transport.ErrRequestedRangeNotSatisfiable: StatusRangeNotSatisfiable,
	transport.ErrTeapot:                       StatusTeapot,
	transport.ErrUnprocessableEntity:          StatusUnprocessableEntity,
	transport.ErrLocked:                       StatusLocked,
	transport.ErrFailedDependency:             StatusFailedDependency,
	transport.ErrPreconditionRequired:         StatusPreconditionRequired,
	transport.ErrTooManyRequests:              StatusTooManyRequests,
	transport.ErrUnavailableForLegalReasons:   StatusUnavailableForLegalReasons,
	// Server
	transport.ErrInternalServerError: StatusInternalServerError,
	transport.ErrNotImplemented:      StatusNotImplemented,
	transport.ErrBadGateway:          StatusBadGateway,
	transport.ErrServiceUnavailable:  StatusServiceUnavailable,
	transport.ErrGatewayTimeout:      StatusGatewayTimeout,
	transport.ErrLoopDetected:        StatusLoopDetected,
}
