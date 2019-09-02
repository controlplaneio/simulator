package simulator_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Config(t *testing.T) {
	t.Skip("Need to mock out terraform output")
	t.Parallel()
	cfg, err := simulator.Config(noopLogger, fixture("noop-tf-dir"), fixture("valid"), "test")

	assert.Nil(t, err)
	assert.NotNil(t, cfg)
}
