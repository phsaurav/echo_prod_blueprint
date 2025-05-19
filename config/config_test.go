package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test helper functions
func TestHelperFunctions(t *testing.T) {
	t.Run("EnvOrDefault", func(t *testing.T) {
		// Test with existing env var
		os.Setenv("TEST_ENV_VAR", "test-value")
		assert.Equal(t, "test-value", envOrDefault("TEST_ENV_VAR", "default"))

		// Test with non-existent env var
		os.Unsetenv("TEST_ENV_VAR")
		assert.Equal(t, "default", envOrDefault("TEST_ENV_VAR", "default"))
	})

	t.Run("ParseDuration", func(t *testing.T) {
		// Valid duration
		assert.Equal(t, 30*time.Second, parseDuration("30s"))

		// Invalid duration
		assert.Equal(t, time.Duration(0), parseDuration("invalid"))
	})

	t.Run("ParseInt", func(t *testing.T) {
		// Valid int
		assert.Equal(t, 42, parseInt("42"))

		// Invalid int
		assert.Equal(t, 0, parseInt("invalid"))
	})

	t.Run("ParseBool", func(t *testing.T) {
		// Valid bools
		assert.Equal(t, true, parseBool("true"))
		assert.Equal(t, false, parseBool("false"))

		// Invalid bool
		assert.Equal(t, false, parseBool("invalid"))
	})

	t.Run("HasPrefix", func(t *testing.T) {
		assert.True(t, hasPrefix(":8080", ":"))
		assert.False(t, hasPrefix("8080", ":"))
		assert.True(t, hasPrefix("test", "t"))
	})
}
