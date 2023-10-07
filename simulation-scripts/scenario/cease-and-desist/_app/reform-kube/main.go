package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

const (
	defaultLicensingServer string = "https://reform-kube.licensing.com.org.test.dev.io"
	productPassword        string = "access-2-reform-kube-server"
)

var (
	password   string
	trial      bool
	licenseUrl string
	help       bool
)

func init() {
	flag.BoolVar(&trial, "trial", false, "Enable program trial mode")
	flag.StringVar(&licenseUrl, "licenseURL", defaultLicensingServer, "Licensing server URL")
	flag.StringVar(&password, "password", "", "Licensing server Password")
	flag.BoolVar(&help, "h", false, "Display help")

	flag.Parse()

	if help {
		fmt.Println("USAGE: <reform-kube-licensing-server> [<args>]")
		flag.PrintDefaults()
	} else {
		if password == "" {
			fmt.Println("Password not set")
			os.Exit(1)
		} else if password != productPassword {
			fmt.Println("Password incorrect")
			os.Exit(1)
		}
	}
}

func main() {

	if trial {
		fmt.Println("Trial mode enabled")
		fmt.Println("FLAG:", os.Getenv("FLAG"))
		licenseFile, _ := json.Marshal(licenseKey{Key: "bGljZW5zZV9rZXk9dHJpYWwK"})
		err := os.WriteFile("license.json", licenseFile, 0444)
		if err != nil {
			fmt.Println("Error writing license file")
			os.Exit(1)
		}
		os.Exit(0)
	}

	err := licenseCheck()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type licenseKey struct {
	Key string `json:"key"`
}

func licenseCheck() error {
	resp, err := http.Get(licenseUrl)
	if err != nil {
		return fmt.Errorf("error contacting licensing server %w", err)
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error contacting licensing server: HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	var licResp licenseKey
	err = json.Unmarshal(body, &licResp)
	if err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}

	if licResp.Key == "bGljZW5zZV9rZXk9cHJvZHVjdGlvbgo=" {
		shim := exec.Command("kubectl", "label", "pods", "rkls", "license=valid", "license_key=2fc593b894ef1402987d2595487d9763", "-n", "licensing")
		_, err := shim.CombinedOutput()
		if err != nil {
			return fmt.Errorf("unable to activate license: %w", err)
		}
		fmt.Println("Product activation successful")
		return nil
	}

	return fmt.Errorf("invalid license key")

}
