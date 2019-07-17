package util_test

import (
	"github.com/controlplaneio/simulator-standalone/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DetectPublicIP_returns_an_IP(t *testing.T) {
	result, err := util.DetectPublicIP()

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, result, "Result was nil")
	assert.NotEmpty(t, result, "Result was empty")
}
