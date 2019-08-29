package util

import (
	"github.com/glendc/go-external-ip"
	"github.com/pkg/errors"
)

// DetectPublicIP detects your public IP address and returns a pointer to a
// string containing the IP address or any error
func DetectPublicIP() (*string, error) {
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		return nil, errors.Wrap(err, "Error detecting public IP address")
	}

	output := ip.String()
	return &output, nil
}
