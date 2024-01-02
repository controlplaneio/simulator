package tools

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strings"
)

type Executable string

type runner struct {
	Executable Executable
	WorkingDir string
	Env        []string
	Arguments  []string
}

func (c runner) Run(ctx context.Context, output io.Writer) error {
	slog.Info("running", "runner", c)

	//nolint:gosec
	cmd := exec.CommandContext(ctx, string(c.Executable), c.Arguments...)
	cmd.Dir = c.WorkingDir
	cmd.Stdout = output
	cmd.Stderr = output
	cmd.Env = c.Env

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run runner: %w", err)
	}

	return nil
}

func (c runner) LogValue() slog.Value {
	cmd := fmt.Sprintf("%s %s", c.Executable, strings.Join(c.Arguments, " "))
	env := make([]string, 0)

	return slog.GroupValue(
		slog.String("cmd", cmd),
		slog.String("dir", c.WorkingDir),
		slog.String("env", strings.Join(env, ",")),
	)
}
