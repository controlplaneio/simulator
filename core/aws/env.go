package aws

import (
	"fmt"
	"os"
)

// envKeys is a list of environment variables that are used by the AWS SDK.
var envKeys = []string{ //nolint: gochecknoglobals
	"AWS_PROFILE",
	"AWS_REGION",
	"AWS_ACCESS_KEY_ID",
	"AWS_SECRET_ACCESS_KEY",
	"AWS_SESSION_TOKEN",
}

func EnvVars() []string {
	env := make([]string, 0)

	for _, key := range envKeys {
		value := os.Getenv(key)
		if len(value) > 0 {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return env
}

func SharedConfigDir(user string) string {
	if user == "root" {
		return "/root/.aws"
	}

	return fmt.Sprintf("/home/%s/.aws", user)
}
