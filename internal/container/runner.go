package container

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/controlplaneio/simulator/v2/controlplane"
	"github.com/controlplaneio/simulator/v2/controlplane/aws"
	"github.com/controlplaneio/simulator/v2/internal/config"
)

const (
	ownerReadWriteExecute = 0700
)

type Simulator interface {
	Run(ctx context.Context, command []string) error
}

func New(config *config.Config) Simulator {
	return &simulator{
		Config: config,
	}
}

type simulator struct {
	Config *config.Config
}

func (r simulator) Run(ctx context.Context, command []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return errors.Join(errors.New("failed to determine home directory"), err)
	}

	// TODO: work with env var for directory
	localAdminSSHBundleDir := filepath.Join(home, ".simulator/admin")
	localPlayerSSHBundleDir := filepath.Join(home, ".simulator/player")
	localAWSDir := filepath.Join(home, ".aws")

	err = mkdirsIfNotExisting(localAdminSSHBundleDir, localPlayerSSHBundleDir)
	if err != nil {
		return errors.Join(errors.New("failed to create bundle directory"), err)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.Join(errors.New("failed to create docker client"), err)
	}

	mounts := []mount.Mount{
		{
			Type:     mount.TypeBind,
			Source:   localAdminSSHBundleDir,
			Target:   controlplane.AdminSSHBundleDir,
			ReadOnly: false,
		},
		{
			Type:     mount.TypeBind,
			Source:   localPlayerSSHBundleDir,
			Target:   controlplane.PlayerSSHBundleDir,
			ReadOnly: false,
		},
		{
			Type:   mount.TypeBind,
			Source: localAWSDir,
			Target: controlplane.AWSDir,
		},
	}

	if r.Config.Cli.Dev {
		mounts = append(mounts, []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: filepath.Join(r.Config.BaseDir, controlplane.Ansible),
				Target: controlplane.AnsibleDir,
			},
			{
				Type:   mount.TypeBind,
				Source: filepath.Join(r.Config.BaseDir, controlplane.Packer),
				Target: controlplane.PackerTemplateDir,
			},
			{
				Type:     mount.TypeBind,
				Source:   filepath.Join(r.Config.BaseDir, controlplane.Terraform),
				Target:   controlplane.TerraformDir,
				ReadOnly: false,
			},
		}...)
	}

	containerConfig := &container.Config{
		Image:        r.Config.Container.Image,
		Env:          aws.Env,
		Cmd:          command,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}

	if r.Config.Container.Rootless {
		// map to host user for directory access
		containerConfig.User = "0:0"
	}

	cont, err := cli.ContainerCreate(ctx,
		containerConfig,
		&container.HostConfig{
			Mounts: mounts,
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		"",
	)
	if err != nil {
		return errors.Join(errors.New("failed to create container"), err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) //nolint:gomnd
		defer cancel()

		err = cli.ContainerStop(ctx, cont.ID, container.StopOptions{})
		if err != nil {
			slog.Warn("failed to stop container", "id", cont.ID, "err", err)
		}

		err = cli.ContainerRemove(ctx, cont.ID, types.ContainerRemoveOptions{})
		if err != nil {
			slog.Warn("failed to remove container", "id", cont.ID, "err", err)
		}
	}()

	hijack, err := cli.ContainerAttach(ctx, cont.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return errors.Join(errors.New("failed to attach to container"), err)
	}

	err = cli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return errors.Join(errors.New("failed to start container"), err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, _ = io.Copy(os.Stdout, hijack.Reader)
		defer wg.Done()
	}()

	wg.Wait()

	return nil
}

func mkdirsIfNotExisting(dirs ...string) error {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(dir, ownerReadWriteExecute)
			if err != nil {
				return errors.Join(errors.New("failed to create directory"), err)
			}
		}
	}
	return nil
}
