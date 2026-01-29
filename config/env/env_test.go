package env

import (
	"encoding/base64"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("load valid .env file", func(t *testing.T) {
		t.Cleanup(func() {
			os.Unsetenv("TEST_VAR")
			os.Unsetenv("ANOTHER_VAR")
		})
		// Create temp env file
		envContent := []byte("TEST_VAR=hello\nANOTHER_VAR=world")
		file := tmpDir + "/.env"
		err := os.WriteFile(file, envContent, 0644)
		require.NoError(t, err)

		// Call New
		err = New(file)
		require.NoError(t, err)

		// Check env vars
		require.Equal(t, "hello", os.Getenv("TEST_VAR"))
		require.Equal(t, "world", os.Getenv("ANOTHER_VAR"))
	})

	t.Run("load multiple .env files", func(t *testing.T) {
		t.Cleanup(func() {
			os.Unsetenv("VAR1")
			os.Unsetenv("VAR2")
			os.Unsetenv("VAR3")
		})

		// Create temp env file1
		file1 := tmpDir + "/file1.env"
		envContent1 := []byte("VAR1=one\nVAR2=two")
		require.NoError(t, os.WriteFile(file1, envContent1, 0644))

		// Create temp env file2
		file2 := tmpDir + "/file2.env"
		envContent2 := []byte("VAR2=override\nVAR3=three")
		require.NoError(t, os.WriteFile(file2, envContent2, 0644))

		// Call New
		err := New(file1, file2)
		require.NoError(t, err)

		// Check env vars
		require.Equal(t, "one", os.Getenv("VAR1"))   // from file1
		require.Equal(t, "two", os.Getenv("VAR2"))   // from file1 (file2 not override file1)
		require.Equal(t, "three", os.Getenv("VAR3")) // from file2
	})

	t.Run(".env file not found", func(t *testing.T) {
		err := New("file/does/not/exist.env")
		require.Error(t, err)
	})
}

// bool

func TestBool(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BOOL") })

		os.Unsetenv("MY_BOOL")
		b, err := Bool("MY_BOOL")
		require.False(t, b)
		require.ErrorContains(t, err, "variable MY_BOOL is not set")
	})

	t.Run("set true", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BOOL") })

		os.Setenv("MY_BOOL", "true")
		b, err := Bool("MY_BOOL")
		require.NoError(t, err)
		require.True(t, b)
	})

	t.Run("set false", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BOOL") })

		os.Setenv("MY_BOOL", "false")
		b, err := Bool("MY_BOOL")
		require.NoError(t, err)
		require.False(t, b)
	})

	t.Run("set invalid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BOOL") })

		os.Setenv("MY_BOOL", "invalid")
		b, err := Bool("MY_BOOL")
		require.False(t, b)
		require.ErrorContains(t, err, "failed to parse MY_BOOL bool value")
	})
}

func TestBoolDefault(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BOOL") })

		os.Unsetenv("MY_BOOL")
		require.True(t, BoolDefault("MY_BOOL", true))
		require.False(t, BoolDefault("MY_BOOL", false))
	})

	t.Run("set true", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BOOL") })

		os.Setenv("MY_BOOL", "true")
		require.True(t, BoolDefault("MY_BOOL", false))
	})

	t.Run("set false", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BOOL") })

		os.Setenv("MY_BOOL", "false")
		require.False(t, BoolDefault("MY_BOOL", true))
	})

	t.Run("set invalid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BOOL") })

		os.Setenv("MY_BOOL", "invalid")
		require.True(t, BoolDefault("MY_BOOL", true))
		require.False(t, BoolDefault("MY_BOOL", false))
	})
}

// int

func TestInt(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_INT") })

		os.Unsetenv("MY_INT")
		i, err := Int("MY_INT")
		require.Equal(t, 0, i)
		require.ErrorContains(t, err, "variable MY_INT is not set")
	})

	t.Run("set valid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_INT") })

		os.Setenv("MY_INT", "31337")
		i, err := Int("MY_INT")
		require.NoError(t, err)
		require.Equal(t, 31337, i)
	})

	t.Run("set invalid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_INT") })

		os.Setenv("MY_INT", "notanumber")
		i, err := Int("MY_INT")
		require.Equal(t, 0, i)
		require.ErrorContains(t, err, "failed to parse MY_INT int value")
	})
}

func TestIntDefault(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_INT") })

		os.Unsetenv("MY_INT")
		require.Equal(t, 31337, IntDefault("MY_INT", 31337))
		require.Equal(t, 0, IntDefault("MY_INT", 0))
	})

	t.Run("set valid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_INT") })

		os.Setenv("MY_INT", "31337")
		require.Equal(t, 31337, IntDefault("MY_INT", 0))
		require.Equal(t, 31337, IntDefault("MY_INT", 123))
	})

	t.Run("set invalid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_INT") })

		os.Setenv("MY_INT", "invalid")
		require.Equal(t, 31337, IntDefault("MY_INT", 31337))
		require.Equal(t, 0, IntDefault("MY_INT", 0))
	})
}

// string

func TestStr(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_STR") })

		os.Unsetenv("MY_STR")
		s, err := Str("MY_STR")
		require.Equal(t, "", s)
		require.ErrorContains(t, err, "variable MY_STR is not set")
	})

	t.Run("set", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_STR") })

		os.Setenv("MY_STR", "hello world")
		s, err := Str("MY_STR")
		require.NoError(t, err)
		require.Equal(t, "hello world", s)
	})
}

func TestStrDefault(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_STR") })

		os.Unsetenv("MY_STR")
		require.Equal(t, "default", StrDefault("MY_STR", "default"))
		require.Equal(t, "", StrDefault("MY_STR", ""))
	})

	t.Run("set", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_STR") })

		os.Setenv("MY_STR", "custom value")
		require.Equal(t, "custom value", StrDefault("MY_STR", "default"))
	})
}

// duration

func TestDur(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_DUR") })

		os.Unsetenv("MY_DUR")
		d, err := Dur("MY_DUR")
		require.Equal(t, time.Duration(0), d)
		require.ErrorContains(t, err, "variable MY_DUR is not set")
	})

	t.Run("set valid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_DUR") })

		os.Setenv("MY_DUR", "2h30m")
		d, err := Dur("MY_DUR")
		require.NoError(t, err)
		require.Equal(t, 2*time.Hour+30*time.Minute, d)
	})

	t.Run("set invalid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_DUR") })

		os.Setenv("MY_DUR", "invalid")
		d, err := Dur("MY_DUR")
		require.Equal(t, time.Duration(0), d)
		require.ErrorContains(t, err, "failed to parse MY_DUR duration value")
	})
}

func TestDurDefault(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_DUR") })

		os.Unsetenv("MY_DUR")
		def := 5 * time.Minute
		require.Equal(t, def, DurDefault("MY_DUR", def))
	})

	t.Run("set valid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_DUR") })

		os.Setenv("MY_DUR", "1h15m")
		require.Equal(t, 1*time.Hour+15*time.Minute, DurDefault("MY_DUR", 10*time.Minute))
	})

	t.Run("set invalid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_DUR") })

		os.Setenv("MY_DUR", "invalid")
		def := 42 * time.Second
		require.Equal(t, def, DurDefault("MY_DUR", def))
	})
}

// bytes (hex)

func TestBytesHex(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BYTES") })

		os.Unsetenv("MY_BYTES")
		b, err := BytesHex("MY_BYTES")
		require.Nil(t, b)
		require.ErrorContains(t, err, "variable MY_BYTES is not set")
	})

	t.Run("set valid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BYTES") })

		value := "hello world"
		os.Setenv("MY_BYTES", hex.EncodeToString([]byte(value)))

		b, err := BytesHex("MY_BYTES")
		require.NoError(t, err)
		require.Equal(t, []byte(value), b)
	})

	t.Run("set invalid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BYTES") })

		os.Setenv("MY_BYTES", "invalid")
		b, err := BytesHex("MY_BYTES")
		require.Nil(t, b)
		require.ErrorContains(t, err, "failed to decode MY_BYTES hex value")
	})
}

func TestBytesHexDefault(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BYTES") })

		def := []byte("hello world")
		os.Unsetenv("MY_BYTES")
		require.Equal(t, def, BytesHexDefault("MY_BYTES", def))
	})

	t.Run("set valid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BYTES") })

		value := "hello world"
		os.Setenv("MY_BYTES", hex.EncodeToString([]byte(value)))

		def := []byte{0x00}
		require.Equal(t, []byte(value), BytesHexDefault("MY_BYTES", def))
	})

	t.Run("set invalid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BYTES") })

		os.Setenv("MY_BYTES", "invalid")
		def := []byte{0xAA, 0xBB}
		require.Equal(t, def, BytesHexDefault("MY_BYTES", def))
	})
}

// bytes (base64)

func TestBytesB64(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BYTES") })

		os.Unsetenv("MY_BYTES")
		b, err := BytesB64("MY_BYTES")
		require.Nil(t, b)
		require.ErrorContains(t, err, "variable MY_BYTES is not set")
	})

	t.Run("set valid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BYTES") })

		value := "hello world"
		os.Setenv("MY_BYTES", base64.StdEncoding.EncodeToString([]byte(value)))

		b, err := BytesB64("MY_BYTES")
		require.NoError(t, err)
		require.Equal(t, []byte(value), b)
	})

	t.Run("set invalid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BYTES") })

		os.Setenv("MY_BYTES", "invalid")
		b, err := BytesB64("MY_BYTES")
		require.Nil(t, b)
		require.ErrorContains(t, err, "failed to decode MY_BYTES base64 value")
	})
}

func TestBytesB64Default(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BYTES") })

		def := []byte("hello world")
		os.Unsetenv("MY_BYTES")
		require.Equal(t, def, BytesB64Default("MY_BYTES", def))
	})

	t.Run("set valid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BYTES") })

		value := "hello world"
		os.Setenv("MY_BYTES", base64.StdEncoding.EncodeToString([]byte(value)))

		def := []byte{0x00}
		require.Equal(t, []byte(value), BytesB64Default("MY_BYTES", def))
	})

	t.Run("set invalid", func(t *testing.T) {
		t.Cleanup(func() { os.Unsetenv("MY_BYTES") })

		os.Setenv("MY_BYTES", "invalid")
		def := []byte{0xAA, 0xBB}
		require.Equal(t, def, BytesB64Default("MY_BYTES", def))
	})
}

