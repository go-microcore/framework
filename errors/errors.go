package errors

import (
	"errors"
	"fmt"

	_ "go.microcore.dev/framework"
)

var (
	// Client errors — issues caused by the client's request

	// The request was invalid or malformed, meaning the server could not
	// understand it due to incorrect syntax, missing fields, or invalid
	// parameters.
	ErrBadRequest error = errors.New("bad_request")

	// The client is not authenticated. This occurs when a request
	// requires authentication but the client did not provide valid
	// credentials or token.
	ErrUnauthorized error = errors.New("unauthorized")

	// Payment is required to access the resource. This is used when
	// access to the requested resource requires a payment, license,
	// or subscription.
	ErrPaymentRequired error = errors.New("payment_required")

	// The client does not have permission to access the resource, even
	// if authenticated. This happens when the user lacks the necessary
	// roles or privileges.
	ErrForbidden error = errors.New("forbidden")

	// The requested resource could not be found on the server. This
	// indicates that the URL is valid but the resource does not exist
	// or has been removed.
	ErrNotFound error = errors.New("not_found")

	// The request method is not allowed for the specified resource.
	// This occurs when, for example, POST is used on a read-only endpoint.
	ErrMethodNotAllowed error = errors.New("method_not_allowed")

	// There is a conflict with the current state of the resource,
	// such as attempting to create a resource that already exists
	// or modify one in an inconsistent state.
	ErrConflict error = errors.New("conflict")

	// The requested resource is no longer available permanently.
	// This is used when a resource has been intentionally removed and
	// will not be available again.
	ErrGone error = errors.New("gone")

	// The request entity is too large for the server to process.
	// This happens when the client sends a payload exceeding server limits.
	ErrRequestEntityTooLarge error = errors.New("request_entity_too_large")

	// The media type of the request is not supported by the server.
	// For example, sending XML to an endpoint that only accepts JSON.
	ErrUnsupportedMediaType error = errors.New("unsupported_media_type")

	// The requested range of data cannot be satisfied. This occurs when
	// a client requests a portion of a resource that is invalid or unavailable.
	ErrRequestedRangeNotSatisfiable error = errors.New("requested_range_not_satisfiable")

	// The server refuses to brew coffee as a teapot. This is a humorous
	// RFC-defined response, rarely used in real applications.
	ErrTeapot error = errors.New("teapot")

	// The server cannot process the request due to semantic errors.
	// The request is well-formed but contains logical errors that prevent
	// it from being processed.
	ErrUnprocessableEntity error = errors.New("unprocessable_entity")

	// The resource is currently locked and cannot be accessed. This usually
	// occurs when another process holds an exclusive lock on the resource.
	ErrLocked error = errors.New("locked")

	// A dependent operation failed, preventing the current operation from
	// completing successfully. Often used in systems with multi-step
	// transactions or workflows.
	ErrFailedDependency error = errors.New("failed_dependency")

	// The request requires a precondition to be met, such as an expected
	// version or ETag. This prevents operations that would conflict with
	// concurrent changes.
	ErrPreconditionRequired error = errors.New("precondition_required")

	// The client has sent too many requests in a given time frame.
	// This is used to implement rate limiting or prevent abuse.
	ErrTooManyRequests error = errors.New("too_many_requests")

	// The resource is unavailable due to legal or regulatory reasons.
	// Access is blocked by government or court orders.
	ErrUnavailableForLegalReasons error = errors.New("unavailable_for_legal_reasons")

	// Server errors — issues caused by server failure or unavailability

	// The server encountered an unexpected condition that prevented it from
	// fulfilling the request. Usually indicates a bug or misconfiguration.
	ErrInternalServerError error = errors.New("internal_server_error")

	// The server does not support the requested functionality. This
	// indicates that the requested method or endpoint has not been
	// implemented.
	ErrNotImplemented error = errors.New("not_implemented")

	// The server received an invalid response from an upstream server
	// or gateway while acting as a proxy. Usually indicates issues in
	// downstream services.
	ErrBadGateway error = errors.New("bad_gateway")

	// The server is temporarily unable to handle the request due to
	// overload, maintenance, or downtime. The condition may resolve
	// after some time.
	ErrServiceUnavailable error = errors.New("service_unavailable")

	// The server did not receive a timely response from an upstream
	// server or dependency while acting as a gateway or proxy.
	ErrGatewayTimeout error = errors.New("gateway_timeout")

	// A loop was detected in processing the request, often due to
	// misconfigured routing, redirects, or recursive dependencies.
	ErrLoopDetected error = errors.New("loop_detected")
)

func New(base error, message string) error {
	return fmt.Errorf("%w:%s", base, message)
}
