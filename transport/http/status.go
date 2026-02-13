package http

import (
	_ "go.microcore.dev/framework"
)

// StatusCode represents an HTTP response status code.
//
// It is defined as a distinct type to improve type safety and readability
// when working with HTTP responses across the framework.
type StatusCode int

// Informational responses (1xx)
//
// These status codes indicate that the request was received and understood,
// and the process is continuing.
const (
	// StatusContinue indicates that the client should continue the request.
	StatusContinue StatusCode = 100

	// StatusSwitchingProtocols indicates that the server is switching protocols.
	StatusSwitchingProtocols StatusCode = 101

	// StatusProcessing indicates that the server has received and is processing the request.
	StatusProcessing StatusCode = 102

	// StatusEarlyHints indicates that the server is sending preliminary response headers.
	StatusEarlyHints StatusCode = 103
)

// Successful responses (2xx)
//
// These status codes indicate that the client's request was successfully
// received, understood, and accepted.
const (
	// StatusOK indicates that the request has succeeded.
	StatusOK StatusCode = 200

	// StatusCreated indicates that a new resource has been successfully created.
	StatusCreated StatusCode = 201

	// StatusAccepted indicates that the request has been accepted for processing.
	StatusAccepted StatusCode = 202

	// StatusNonAuthoritativeInfo indicates that returned metadata is from a local or third-party copy.
	StatusNonAuthoritativeInfo StatusCode = 203

	// StatusNoContent indicates success with no response body.
	StatusNoContent StatusCode = 204

	// StatusResetContent indicates that the client should reset the document view.
	StatusResetContent StatusCode = 205

	// StatusPartialContent indicates that only part of the resource is being returned.
	StatusPartialContent StatusCode = 206

	// StatusMultiStatus provides multiple status values for different operations.
	StatusMultiStatus StatusCode = 207

	// StatusAlreadyReported indicates that members have already been enumerated.
	StatusAlreadyReported StatusCode = 208

	// StatusIMUsed indicates that instance manipulations were applied to the response.
	StatusIMUsed StatusCode = 226
)

// Redirection messages (3xx)
//
// These status codes indicate that further action needs to be taken by the
// client to complete the request.
const (
	// StatusMultipleChoices indicates multiple resource options for the client.
	StatusMultipleChoices StatusCode = 300

	// StatusMovedPermanently indicates that the resource has been permanently moved.
	StatusMovedPermanently StatusCode = 301

	// StatusFound indicates that the resource resides temporarily under a different URI.
	StatusFound StatusCode = 302

	// StatusSeeOther indicates that the client should retrieve the resource at another URI.
	StatusSeeOther StatusCode = 303

	// StatusNotModified indicates that the cached resource has not changed.
	StatusNotModified StatusCode = 304

	// StatusUseProxy indicates that the request must be accessed through a proxy (deprecated).
	StatusUseProxy StatusCode = 305

	// StatusTemporaryRedirect indicates a temporary redirect preserving the request method.
	StatusTemporaryRedirect StatusCode = 307

	// StatusPermanentRedirect indicates a permanent redirect preserving the request method.
	StatusPermanentRedirect StatusCode = 308
)

// Client error responses (4xx)
//
// These status codes indicate that the client seems to have made an error
// in the request.
const (
	// StatusBadRequest indicates that the server cannot process the request due to client error.
	StatusBadRequest StatusCode = 400

	// StatusUnauthorized indicates that authentication is required.
	StatusUnauthorized StatusCode = 401

	// StatusPaymentRequired is reserved for future use.
	StatusPaymentRequired StatusCode = 402

	// StatusForbidden indicates that the server refuses to authorize the request.
	StatusForbidden StatusCode = 403

	// StatusNotFound indicates that the requested resource was not found.
	StatusNotFound StatusCode = 404

	// StatusMethodNotAllowed indicates that the HTTP method is not allowed for the resource.
	StatusMethodNotAllowed StatusCode = 405

	// StatusNotAcceptable indicates that no acceptable representation is available.
	StatusNotAcceptable StatusCode = 406

	// StatusProxyAuthRequired indicates that proxy authentication is required.
	StatusProxyAuthRequired StatusCode = 407

	// StatusRequestTimeout indicates that the server timed out waiting for the request.
	StatusRequestTimeout StatusCode = 408

	// StatusConflict indicates a conflict with the current resource state.
	StatusConflict StatusCode = 409

	// StatusGone indicates that the resource is permanently unavailable.
	StatusGone StatusCode = 410

	// StatusLengthRequired indicates that the request must specify Content-Length.
	StatusLengthRequired StatusCode = 411

	// StatusPreconditionFailed indicates that a request precondition failed.
	StatusPreconditionFailed StatusCode = 412

	// StatusPayloadTooLarge indicates that the request payload is too large.
	StatusPayloadTooLarge StatusCode = 413

	// StatusURITooLong indicates that the request URI is too long.
	StatusURITooLong StatusCode = 414

	// StatusUnsupportedMediaType indicates that the media type is not supported.
	StatusUnsupportedMediaType StatusCode = 415

	// StatusRangeNotSatisfiable indicates that the requested range cannot be fulfilled.
	StatusRangeNotSatisfiable StatusCode = 416

	// StatusExpectationFailed indicates that an Expect header requirement could not be met.
	StatusExpectationFailed StatusCode = 417

	// StatusTeapot indicates that the server refuses to brew coffee because it is a teapot.
	StatusTeapot StatusCode = 418

	// StatusMisdirectedRequest indicates that the request was sent to the wrong server.
	StatusMisdirectedRequest StatusCode = 421

	// StatusUnprocessableEntity indicates that the request is semantically invalid.
	StatusUnprocessableEntity StatusCode = 422

	// StatusLocked indicates that the resource is locked.
	StatusLocked StatusCode = 423

	// StatusFailedDependency indicates that the request failed due to another failed request.
	StatusFailedDependency StatusCode = 424

	// StatusTooEarly indicates that the request may be replayed.
	StatusTooEarly StatusCode = 425

	// StatusUpgradeRequired indicates that the client must switch to a different protocol.
	StatusUpgradeRequired StatusCode = 426

	// StatusPreconditionRequired indicates that the request must be conditional.
	StatusPreconditionRequired StatusCode = 428

	// StatusTooManyRequests indicates that the client has sent too many requests.
	StatusTooManyRequests StatusCode = 429

	// StatusRequestHeaderFieldsTooLarge indicates that request headers are too large.
	StatusRequestHeaderFieldsTooLarge StatusCode = 431

	// StatusUnavailableForLegalReasons indicates that the resource is blocked for legal reasons.
	StatusUnavailableForLegalReasons StatusCode = 451
)

// Server error responses (5xx)
//
// These status codes indicate that the server failed to fulfill a valid request.
const (
	// StatusInternalServerError indicates that the server encountered an unexpected condition.
	StatusInternalServerError StatusCode = 500

	// StatusNotImplemented indicates that the server does not support the requested functionality.
	StatusNotImplemented StatusCode = 501

	// StatusBadGateway indicates that the server received an invalid response from an upstream server.
	StatusBadGateway StatusCode = 502

	// StatusServiceUnavailable indicates that the server is temporarily unable to handle the request.
	StatusServiceUnavailable StatusCode = 503

	// StatusGatewayTimeout indicates that the upstream server did not respond in time.
	StatusGatewayTimeout StatusCode = 504

	// StatusHTTPVersionNotSupported indicates that the HTTP version is not supported.
	StatusHTTPVersionNotSupported StatusCode = 505

	// StatusVariantAlsoNegotiates indicates a configuration error in content negotiation.
	StatusVariantAlsoNegotiates StatusCode = 506

	// StatusInsufficientStorage indicates that the server cannot store the representation.
	StatusInsufficientStorage StatusCode = 507

	// StatusLoopDetected indicates that an infinite loop was detected.
	StatusLoopDetected StatusCode = 508

	// StatusNotExtended indicates that further extensions are required to fulfill the request.
	StatusNotExtended StatusCode = 510

	// StatusNetworkAuthenticationRequired indicates that network authentication is required.
	StatusNetworkAuthenticationRequired StatusCode = 511
)
