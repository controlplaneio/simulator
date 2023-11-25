package tools

import (
	"context"
	"errors"
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
	Arguments  []string
}

func (c runner) Run(ctx context.Context, output io.Writer) error {
	slog.Info("running", "runner", c)

	//nolint:gosec
	cmd := exec.CommandContext(ctx, string(c.Executable), c.Arguments...)
	cmd.Dir = c.WorkingDir
	cmd.Stdout = output
	cmd.Stderr = output

	err := cmd.Run()
	if err != nil {
		return errors.Join(errors.New("failed to run runner"), err)
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
