package commands

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

type Executable string

type Runnable interface {
	Run(ctx context.Context, output ...io.Writer) error
}

type command struct {
	Executable  Executable
	WorkingDir  string
	Environment []string
	Arguments   []string
}

func (c command) Run(ctx context.Context, output ...io.Writer) error {
	slog.Info("running", "command", c)

	// Default to writing to stdout unless an alternative writer is provided
	var writer io.Writer
	if len(output) == 0 {
		writer = os.Stdout
	} else {
		writer = output[0]
	}

	cmd := exec.CommandContext(ctx, string(c.Executable), c.Arguments...)
	cmd.Dir = c.WorkingDir
	cmd.Env = c.Environment
	cmd.Stdout = writer
	cmd.Stderr = writer

	// TODO: Ensure ctrl-c stops the command
	err := cmd.Run()
	if err != nil {
		slog.Error("failed to run", "command", c)
	}

	return err
}

func (c command) LogValue() slog.Value {
	cmd := fmt.Sprintf("%s %s", c.Executable, strings.Join(c.Arguments, " "))
	var env []string

	// Only log env keys, not values
	for _, value := range c.Environment {
		env = append(env, value[:strings.Index(value, "=")])
	}

	return slog.GroupValue(
		slog.String("cmd", cmd),
		slog.String("dir", c.WorkingDir),
		slog.String("env", strings.Join(env, ",")),
	)
}
