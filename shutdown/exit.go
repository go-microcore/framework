package shutdown // import "go.microcore.dev/framework/shutdown"

import (
	"errors"
	"fmt"
)

// ExitReason represents an error that carries an exit code for program termination.
// It wraps an underlying error, allowing both the error message and the exit code
// to be propagated together. Use NewExitReason(code, err) to create an ExitReason.

type ExitReason struct {
	Code int
	Err  error
}

func (e *ExitReason) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("exit code: %d", e.Code)
	}
	return e.Err.Error()
}

func (e *ExitReason) Unwrap() error {
	return e.Err
}

// NewExitReason creates a new program termination reason with the specified exit code.
//
// ExitReason represents a reason for program termination. It can wrap underlying errors (if provided)
// and associates them with a machine-readable exit code. This is useful for CLI applications or services
// where the process exit code matters (for orchestration, CI/CD, or supervisors).
//
// The underlying errors `errs` can be omitted if you only want to signal an exit code
// without providing actual error messages.
//
// Example:
//
//	// Exit with a software error and a descriptive message
//	return NewExitReason(ExitSoftware, fmt.Errorf("failed to initialize database"))
//
//	// Exit with just a code, no underlying error
//	return NewExitReason(ExitUnavailable)
func NewExitReason(code int, errs ...error) error {
	return &ExitReason{
		Code: code,
		Err:  errors.Join(errs...),
	}
}

// ParseExitReason extracts the exit code and error presence from a wrapped ExitReason.
//
// This function recursively unwraps the provided error `e` looking for ExitReason instances.
// It returns two values:
//   - code: the most relevant exit code found (defaults to ExitGeneralError if none found)
//   - hasErr: true if there is an actual underlying error, false if the ExitReason was just signaling a code
//
// ParseExitReason is intended for CLI applications or services that need to translate
// error chains into appropriate exit codes for process termination.
//
// Example usage:
//
//	// Simple case: just an ExitReason with a code
//	reason := NewExitReason(ExitUnavailable)
//	code, hasErr := ParseExitReason(reason)
//	fmt.Println(code, hasErr) // Output: ExitUnavailable, false
//	fmt.Println(reason.Error()) // Output: exit code: 69
//
//	// ExitReason wrapping an underlying error
//	reason = NewExitReason(ExitSoftware, fmt.Errorf("failed to initialize database"))
//	code, hasErr = ParseExitReason(reason)
//	fmt.Println(code, hasErr) // Output: ExitSoftware, true
//	fmt.Println(reason.Error()) // Output: failed to initialize database
//
//	// Nested ExitReasons: inner ExitReason can be preserved
//	reason1 := NewExitReason(ExitUnavailable, fmt.Errorf("reason1"))
//	reason2 := NewExitReason(ExitSoftware, fmt.Errorf("reason2: %w", reason1))
//	code, hasErr = ParseExitReason(reason2)
//	fmt.Println(code, hasErr) // Output: ExitUnavailable, true
//	fmt.Println(reason2.Error()) // Output: reason2: reason1
//
//	// No ExitReason in chain: defaults to general error
//	err := fmt.Errorf("just a normal error")
//	code, hasErr = ParseExitReason(err)
//	fmt.Println(code, hasErr) // Output: ExitGeneralError, true
//	fmt.Println(err.Error()) // Output: just a normal error
//
// Use this function whenever you need a single exit code for the process while still
// capturing whether there was an actual error message to log or display.
func ParseExitReason(e error) (int, bool) {
	code := ExitGeneralError // default exit code
	err := true              // default error state
	cur := e
	for cur != nil {
		ee := &ExitReason{}
		if errors.As(cur, &ee) {
			code = ee.Code
			err = ee.Err != nil
			cur = ee.Err
			continue
		}
		cur = nil
	}
	return code, err
}

// These constants are intended to be used as a stable contract between
// Go applications and their execution environment (OS, Docker, Kubernetes,
// systemd, CI/CD pipelines, supervisors).
//
// Exit codes MUST NOT be used as a logging mechanism.
// They are a machine-readable signal describing WHY the process terminated.
//
// Conventions:
//   - 0      → successful termination
//   - 1–63   → generic / application-defined errors
//   - 64–78  → sysexits (BSD / POSIX de-facto standard)
//   - 128+N  → terminated by Unix signal N

const (
	// ExitOK indicates successful termination.
	//
	// Use when:
	//   - the application completed its work successfully
	//   - a server shut down gracefully after receiving SIGTERM
	//   - a CLI command or seed finished without errors
	//
	// MUST be used for normal, expected shutdown.
	ExitOK = 0

	// ExitGeneralError indicates an unspecified failure.
	//
	// Use when:
	//   - an error occurred but does not fit a more specific category
	//   - acting as a fallback error code
	//
	// Avoid using this when a more precise exit code exists.
	ExitGeneralError = 1

	// ExitUsage indicates incorrect command usage.
	//
	// Use when:
	//   - CLI arguments are invalid
	//   - required flags are missing
	//   - incompatible flags are provided
	//
	// Typical for CLI tools and admin commands.
	ExitUsage = 64

	// ExitDataError indicates invalid input data.
	//
	// Use when:
	//   - input data is malformed
	//   - JSON / YAML / CSV parsing fails due to invalid content
	//   - semantic validation fails
	ExitDataError = 65

	// ExitNoInput indicates a missing required input.
	//
	// Use when:
	//   - a required file does not exist
	//   - stdin or expected input source is unavailable
	ExitNoInput = 66

	// ExitUnavailable indicates that a required external service is unavailable.
	//
	// Use when:
	//   - database cannot be reached
	//   - Redis / Kafka / external API is down
	//   - network dependency is unreachable at startup
	//
	// This usually triggers restart in orchestration systems.
	ExitUnavailable = 69

	// ExitSoftware indicates an internal software error.
	//
	// Use when:
	//   - an invariant is violated
	//   - unexpected state is reached
	//   - a bug is detected but panic is not used
	ExitSoftware = 70

	// ExitOSError indicates an operating system error.
	//
	// Use when:
	//   - syscall failures occur
	//   - OS-level resources cannot be accessed
	ExitOSError = 71

	// ExitIOError indicates a low-level I/O failure.
	//
	// Use when:
	//   - disk read/write fails
	//   - socket I/O fails unexpectedly
	ExitIOError = 74

	// ExitTempFail indicates a temporary failure.
	//
	// Use when:
	//   - the operation can be retried
	//   - transient network issues occur
	//   - rate limits are hit
	//
	// Supervisors may retry automatically.
	ExitTempFail = 75

	// ExitNoPermission indicates insufficient permissions.
	//
	// Use when:
	//   - access to secrets is denied
	//   - filesystem permissions are insufficient
	//   - security constraints prevent startup
	ExitNoPermission = 77

	// ExitConfigError indicates an invalid configuration.
	//
	// Use when:
	//   - required environment variables are missing
	//   - configuration files are invalid or malformed
	//   - configuration values fail validation
	//
	// This is one of the most important exit codes for production systems.
	ExitConfigError = 78

	// ExitPanic indicates the application terminated due to a panic.
	//
	// Use when:
	//   - a panic was recovered in main()
	//   - an unrecoverable programming error occurred
	//
	// Panics SHOULD be treated differently from normal errors.
	ExitPanic = 10

	// ExitShutdownError indicates failure during graceful shutdown.
	//
	// Use when:
	//   - shutdown handlers exceed timeout
	//   - resources fail to close cleanly
	//
	// The application attempted to shut down gracefully but failed.
	ExitShutdownError = 20

	// ExitSignalBase is the base exit code for Unix signals.
	//
	// Actual exit code is calculated as:
	//   128 + signal number
	//
	// Examples:
	//   SIGINT  (2)  → 130
	//   SIGTERM (15) → 143
	//   SIGKILL (9)  → 137
	//
	// This value MUST NOT be returned manually.
	ExitSignalBase = 128
)
