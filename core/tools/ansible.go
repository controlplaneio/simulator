package tools

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/controlplaneio/simulator/v2/internal/docker"
)

const (
	AnsiblePlaybookExecutable Executable = "ansible-playbook"
)

type ScenarioManager interface {
	Install(ctx context.Context, id string) error
	Uninstall(ctx context.Context, id string) error
}

type AnsiblePlaybook struct {
	WorkingDir  string
	PlaybookDir string
	Env         []string
	StdOut      io.Writer
	StdErr      io.Writer
}

func (p AnsiblePlaybook) Install(ctx context.Context, id string) error {
	playbook := fmt.Sprintf("%s.yaml", id)

	err := ansiblePlaybookCommand(p.WorkingDir, p.PlaybookDir, p.Env, playbook).Run(ctx, p.StdOut, p.StdErr)
	if err != nil {
		return fmt.Errorf("failed to execute Ansible Playbook: %w", err)
	}

	return nil
}

func (p AnsiblePlaybook) Uninstall(ctx context.Context, id string) error {
	playbook := fmt.Sprintf("%s.yaml", id)

	if err := ansiblePlaybookCommand(p.WorkingDir, p.PlaybookDir, p.Env, playbook, "state=absent").
		Run(ctx, p.StdOut, p.StdErr); err != nil {
		return fmt.Errorf("failed to run Ansible Playbook with state=absent: %w", err)
	}

	return nil
}

func ansiblePlaybookCommand(workingDir, playbookDir string, env []string, playbook string, extraVars ...string) runner {
	args := []string{
		fmt.Sprintf("%s/%s", playbookDir, playbook),
	}

	if len(extraVars) > 0 {
		args = append(args,
			"--extra-vars",
			strings.Join(extraVars, " "),
		)
	}

	return runner{
		Executable: AnsiblePlaybookExecutable,
		WorkingDir: workingDir,
		Arguments:  args,
		Env:        env,
	}
}

type AnsiblePlaybookContainer struct {
	Client *docker.Client
	Config *docker.Config
}

func (p AnsiblePlaybookContainer) Install(ctx context.Context, id string) error {
	config := *p.Config
	config.Cmd = []string{
		"scenario",
		"install",
		id,
	}

	if err := p.Client.Run(ctx, config); err != nil {
		return fmt.Errorf("failed to build ami: %w", err)
	}

	return nil
}

//nolint:varnamelen
func (p AnsiblePlaybookContainer) Uninstall(ctx context.Context, id string) error {
	config := *p.Config
	config.Cmd = []string{
		"scenario",
		"uninstall",
		id,
	}

	if err := p.Client.Run(ctx, config); err != nil {
		return fmt.Errorf("failed to uninstall playbook with id %s: %w", id, err)
	}

	return nil
}
