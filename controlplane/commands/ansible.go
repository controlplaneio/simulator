package commands

import (
	"fmt"
	"strings"
)

const (
	AnsiblePlaybook Executable = "ansible-playbook"
)

func AnsiblePlaybookCommand(workingDir, playbook string, extraVars ...string) Runnable {
	args := []string{
		fmt.Sprintf("%s/%s.yaml", workingDir, playbook),
	}

	if len(extraVars) > 0 {
		args = append(args,
			"--extra-vars",
			fmt.Sprintf("\"%s\"", strings.Join(extraVars, " ")),
		)
	}

	return command{
		Executable:  AnsiblePlaybook,
		WorkingDir:  workingDir,
		Environment: nil,
		Arguments:   args,
	}
}
