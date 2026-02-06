package shutdown

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestManager_NewContext(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	ctx, err := m.NewContext()
	require.NoError(t, err)
	require.NotNil(t, ctx)

	cancelVal := m.ctx.cancel.Load()
	require.NotNil(t, cancelVal)
}

func TestManager_NewContext_Twice(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	_, err := m.NewContext()
	require.NoError(t, err)

	_, err = m.NewContext()
	require.ErrorIs(t, err, ErrContextAlreadyInit)
}

func TestManager_WithContext_CustomParent(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	type ctxKey string
	parent := context.WithValue(context.Background(), ctxKey("key"), "value")
	ctx, err := m.WithContext(parent)
	require.NoError(t, err)
	require.NotNil(t, ctx)

	require.Equal(t, "value", ctx.Value(ctxKey("key")))

	_, err = m.WithContext(context.Background())
	require.ErrorIs(t, err, ErrContextAlreadyInit)
}

func TestManager_WithContext_NilParent(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	//lint:ignore SA1012 WithContext with nil
	_, err := m.WithContext(nil)
	require.ErrorIs(t, err, ErrParentContextNil)
}

func TestManager_Context_ReturnsExistingContext(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	ctx, err := m.NewContext()
	require.NoError(t, err)

	ctx2 := m.Context()
	require.Same(t, ctx, ctx2)
}

func TestManager_Context_ReturnsBackgroundIfNotInitialized(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	ctx := m.Context()
	require.NotNil(t, ctx)
	require.Equal(t, context.Background(), ctx)
}

func TestManager_AddHandler_BeforeShutdown(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	err := m.AddHandler(func(ctx context.Context, code int) error {
		return nil
	})

	require.NoError(t, err)
}

func TestManager_AddHandler_AfterShutdown(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	m.Shutdown(ExitOK)

	time.Sleep(10 * time.Millisecond)

	err := m.AddHandler(func(ctx context.Context, code int) error {
		return nil
	})

	require.ErrorIs(t, err, ErrCannotAddHandlerAfterShutdown)
}

func TestManager_Wait_Success(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	exit := make(chan int, 1)

	go func() {
		exit <- m.Wait()
	}()

	expectedCode := ExitOK
	m.exit <- expectedCode

	select {
	case code := <-exit:
		require.Equal(t, expectedCode, code)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestManager_Wait_Blocks(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	go m.Wait()

	select {
	case <-m.exit:
		t.Fatal("Wait returned without exit code")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestManager_Shutdown_SendsCode(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	var logs []string
	origLogger := logger
	logger = slog.New(slog.NewTextHandler(&mockWriter{logs: &logs}, nil))
	defer func() { logger = origLogger }()

	expectedCode := ExitOK
	m.Shutdown(expectedCode)

	select {
	case code := <-m.code:
		require.Equal(t, expectedCode, code)
	case <-time.After(time.Second):
		t.Fatal("code was not sent to m.code")
	}

	require.Empty(t, logs)
}

func TestManager_Shutdown_ChannelBlocked(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	tmpCode := ExitGeneralError
	m.code <- tmpCode

	var logs []string
	origLogger := logger
	logger = slog.New(slog.NewTextHandler(&mockWriter{logs: &logs}, nil))
	defer func() { logger = origLogger }()

	m.Shutdown(ExitOK)

	select {
	case code := <-m.code:
		require.Equal(t, tmpCode, code)
	default:
		t.Fatal("channel unexpectedly empty")
	}

	require.Len(t, logs, 1)
	require.Contains(t, logs[0], "channel blocked")
}

func TestManager_Exit(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	exit := make(chan int, 1)
	expectedCode := ExitOK

	go func() {
		exit <- m.Exit(expectedCode)
	}()

	select {
	case code := <-exit:
		require.Equal(t, expectedCode, code)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout")
	}
}

func TestManager_SetShutdownTimeout(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	origTimeout := timeout
	defer func() { timeout = origTimeout }()

	newTimeout := 123 * time.Millisecond
	m.SetShutdownTimeout(newTimeout)

	require.Equal(t, newTimeout, timeout)
}

func TestSetDefaultManager_BeforeInit(t *testing.T) {
	custom := newManager()

	err := SetDefaultManager(custom)
	require.NoError(t, err)
}

func TestSetDefaultManager_AfterDefault(t *testing.T) {
	_ = Context() // init default

	err := SetDefaultManager(newManager())
	require.ErrorIs(t, err, ErrManagerAlreadyRunning)
}

func TestManager_ContextCanceledOnShutdown(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	ctx, err := m.NewContext()
	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		close(done)
	}()

	m.Shutdown(ExitOK)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("context was not canceled")
	}
}

func TestManager_HandlerSuccess(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	called := false

	m.AddHandler(func(ctx context.Context, code int) error {
		require.Equal(t, ExitOK, code)
		called = true
		return nil
	})

	m.Shutdown(ExitOK)

	code := <-m.exit
	require.Equal(t, ExitOK, code)
	require.Equal(t, true, called)
}

func TestManager_HandlerError(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	m.AddHandler(func(ctx context.Context, code int) error {
		return errors.New("fail")
	})

	m.Shutdown(ExitOK)

	code := <-m.exit
	require.Equal(t, ExitShutdownError, code)
}

func TestManager_HandlerPanic(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	m.AddHandler(func(ctx context.Context, code int) error {
		panic("boom")
	})

	m.Shutdown(ExitOK)

	code := <-m.exit
	require.Equal(t, ExitShutdownError, code)
}

func TestManager_HandlerTimeout(t *testing.T) {
	t.Parallel()
	m := newManager().(*manager)

	m.SetShutdownTimeout(100 * time.Millisecond)

	m.AddHandler(func(ctx context.Context, code int) error {
		<-ctx.Done()
		return nil
	})

	m.Shutdown(ExitOK)

	code := <-m.exit
	require.Equal(t, ExitShutdownError, code)
}

func TestNewExitReason(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		reason     error
		wantCode   int
		wantHasErr bool
		wantErrMsg string
	}{
		{
			name:       "Only exit code, no error",
			reason:     NewExitReason(ExitUnavailable),
			wantCode:   ExitUnavailable,
			wantHasErr: false,
			wantErrMsg: "exit code: 69",
		},
		{
			name:       "Exit code with single error",
			reason:     NewExitReason(ExitSoftware, errors.New("failed to initialize DB")),
			wantCode:   ExitSoftware,
			wantHasErr: true,
			wantErrMsg: "failed to initialize DB",
		},
		{
			name: "Nested ExitReasons",
			reason: func() error {
				inner := NewExitReason(ExitUnavailable, errors.New("inner error"))
				return NewExitReason(ExitSoftware, fmt.Errorf("outer error: %w", inner))
			}(),
			wantCode:   ExitUnavailable,
			wantHasErr: true,
			wantErrMsg: "outer error: inner error",
		},
		{
			name:       "Normal error, no ExitReason",
			reason:     errors.New("just a normal error"),
			wantCode:   ExitGeneralError,
			wantHasErr: true,
			wantErrMsg: "just a normal error",
		},
		{
			name:       "Multiple errors joined",
			reason:     NewExitReason(ExitSoftware, errors.New("err1"), errors.New("err2")),
			wantCode:   ExitSoftware,
			wantHasErr: true,
			wantErrMsg: "err1\nerr2", // errors.Join joins with newline by default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, hasErr := ParseExitReason(tt.reason)
			if code != tt.wantCode {
				t.Errorf("ParseExitReason() code = %d, want %d", code, tt.wantCode)
			}
			if hasErr != tt.wantHasErr {
				t.Errorf("ParseExitReason() hasErr = %v, want %v", hasErr, tt.wantHasErr)
			}
			if tt.reason.Error() != tt.wantErrMsg {
				t.Errorf("Error() = %q, want %q", tt.reason.Error(), tt.wantErrMsg)
			}
		})
	}
}


type mockWriter struct {
	logs *[]string
}

func (w *mockWriter) Write(p []byte) (int, error) {
	*w.logs = append(*w.logs, string(p))
	return len(p), nil
}
