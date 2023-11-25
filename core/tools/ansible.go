package tools

import (
	"context"
	"errors"
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
	Output      io.Writer
}

func (p AnsiblePlaybook) Install(ctx context.Context, id string) error {
	playbook := fmt.Sprintf("%s.yaml", id)

	if err := ansiblePlaybookCommand(p.WorkingDir, p.PlaybookDir, playbook).Run(ctx, p.Output); err != nil {
		return errors.Join(errors.New("failed to execute Ansible Playbook"), err)
	}

	return nil
}

func (p AnsiblePlaybook) Uninstall(ctx context.Context, id string) error {
	playbook := fmt.Sprintf("%s.yaml", id)

	if err := ansiblePlaybookCommand(p.WorkingDir, p.PlaybookDir, playbook, "state=absent").
		Run(ctx, p.Output); err != nil {
		return errors.Join(errors.New("failed to execute Ansible Playbook"), err)
	}

	return nil
}

func ansiblePlaybookCommand(workingDir, playbookDir, playbook string, extraVars ...string) runner {
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
		return errors.Join(errors.New("failed to build ami"), err)
	}

	return nil
}

func (p AnsiblePlaybookContainer) Uninstall(ctx context.Context, id string) error {
	config := *p.Config
	config.Cmd = []string{
		"scenario",
		"uninstall",
		id,
	}

	if err := p.Client.Run(ctx, config); err != nil {
		return errors.Join(errors.New("failed to build ami"), err)
	}

	return nil
}
