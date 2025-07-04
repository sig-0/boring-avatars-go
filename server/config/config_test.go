package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ValidateConfig(t *testing.T) {
	t.Parallel()

	t.Run("invalid listen address", func(t *testing.T) {
		t.Parallel()

		cfg := DefaultConfig()
		cfg.ListenAddress = "rando-address" // doesn't follow the format

		assert.ErrorIs(t, ValidateConfig(cfg), ErrInvalidListenAddress)
	})

	t.Run("valid configuration", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, ValidateConfig(DefaultConfig()))
	})
}
