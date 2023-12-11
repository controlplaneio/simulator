package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &Client{cli}, nil
}

type Client struct {
	client *client.Client
}

func (c Client) PullImage(ctx context.Context, ref string) error {
	out, err := c.client.ImagePull(ctx, ref, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer func() {
		_ = out.Close()
	}()

	if _, err = io.Copy(os.Stdout, out); err != nil {
		return errors.Join(errors.New("failed to pull image"), err)
	}

	return nil
}

func (c Client) Run(ctx context.Context, conf Config) error {
	// Pre-allocate the 'mounts' slice with the required capacity
	mounts := make([]mount.Mount, len(conf.Mounts))

	for i, m := range conf.Mounts {
		mounts[i] = mount.Mount{
			Type:     mount.TypeBind,
			Source:   m.Source,
			Target:   m.Target,
			ReadOnly: m.ReadOnly,
		}
	}

	containerConfig := &container.Config{
		Image:        conf.Image,
		Env:          conf.Env,
		Cmd:          conf.Cmd,
		Tty:          true,
		AttachStderr: true,
		AttachStdout: true,
	}

	if conf.Rootless {
		// map to host user for directory access
		containerConfig.User = "0:0"
	}

	cont, err := c.client.ContainerCreate(ctx,
		containerConfig,
		&container.HostConfig{
			Mounts: mounts,
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		"",
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	defer func() {
		//nolint:gomnd
		cctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		// TODO: above, do we need to create a new context from the original ctx?
		defer cancel()

		if err := c.client.ContainerStop(cctx, cont.ID, container.StopOptions{}); err != nil {
			slog.Warn("failed to stop container", "id", cont.ID, "err", err)
		}

		if err := c.client.ContainerRemove(cctx, cont.ID, types.ContainerRemoveOptions{}); err != nil {
			slog.Warn("failed to remove container", "id", cont.ID, "err", err)
		}
	}()

	hijack, err := c.client.ContainerAttach(ctx, cont.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return fmt.Errorf("failed to attach to container: %w", err)
	}

	err = c.client.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func() {
		_, err := io.Copy(os.Stdout, hijack.Reader)
		if err != nil {
			slog.Warn("failed to copy container output", "err", err)
		}
		defer waitGroup.Done()
	}()

	waitGroup.Wait()

	return nil
}

type Config struct {
	Image    string
	Env      []string
	Cmd      []string
	Mounts   []MountConfig
	Rootless bool
}

type MountConfig struct {
	Source   string
	Target   string
	ReadOnly bool
}
