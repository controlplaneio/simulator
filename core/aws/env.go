package aws

import (
	"fmt"
	"os"
)

var (
	//nolint: gochecknoglobals
	// envKeys is a list of environment variables that are used by the AWS SDK.
	envKeys = []string{
		"AWS_PROFILE",
		"AWS_REGION",
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_SESSION_TOKEN",
	}
)

func EnvVars() []string {
	env := make([]string, 0)

	for _, key := range envKeys {
		value, ok := os.LookupEnv(key)
		if ok && len(value) > 0 {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return env
}
