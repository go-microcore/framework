package shutdown // import "go.microcore.dev/framework/shutdown"

import "errors"

var (
	// Attempting to create a shutdown context but the root context has already been initialized.
	ErrContextAlreadyInit = errors.New("shutdown context already initialized")

	// Attempting to create a shutdown context with a nil parent.
	ErrParentContextNil = errors.New("parent context is nil")

	// Trying to add a shutdown handler after the shutdown process has started.
	ErrCannotAddHandlerAfterShutdown = errors.New("cannot add handler after shutdown started")

	// SetDefault is called after the default manager is already initialized.
	ErrManagerAlreadyRunning = errors.New("manager already runned")

	// SetDefault is called after shutdown has started or completed.
	ErrCannotCallAfterShutdown = errors.New("cannot call after shutdown started")

	// Unknown state is encountered.
	ErrUnknownState = errors.New("failed due to unknown state")
)
