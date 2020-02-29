package simulator_test

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/controlplaneio/simulator-standalone/pkg/simulator"
)

func Test_Perturb(t *testing.T) {
	os.Setenv("SIMULATOR_SCENARIOS_DIR", fixture("noop-perturb"))
	po := simulator.PerturbOptions{}
	_, err := simulator.Perturb(&po, logrus.New())
	assert.NoError(t, err)
}
