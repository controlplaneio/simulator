package util

import (
	"github.com/jc21/go-external-ip"
	"github.com/pkg/errors"
)

// DetectPublicIP detects your public IP address and returns a pointer to a
// string containing the IP address or any error
func DetectPublicIP() (*string, error) {
	consensus := externalip.DefaultConsensus(nil, nil)
	err := consensus.UseIPProtocol(4)
	if err != nil {
		return nil, errors.Wrap(err, "Error detecting public IP address")
	}

	ip, err := consensus.ExternalIP()
	if err != nil {
		return nil, errors.Wrap(err, "Error detecting public IP address")
	}

	output := ip.String()
	return &output, nil
}
