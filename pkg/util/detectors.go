package util

import (
	"encoding/json"
	"github.com/glendc/go-external-ip"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

// UbuntuCloudImageReleasesURL is a string contianing the URL to find AMI IDs for AWS regions
const UbuntuCloudImageReleasesURL = "https://cloud-images.ubuntu.com/locator/ec2/releasesTable"

// DetectPublicIP detects your public IP address and returns a pointer to a string containing the IP address or any error
func DetectPublicIP() (*string, error) {
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		return nil, err
	}

	output := ip.String()
	return &output, nil
}
