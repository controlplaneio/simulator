package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/controlplaneio/simulator/v2/controlplane/aws"
	"github.com/controlplaneio/simulator/v2/internal/container"
)

var imageCmd = &cobra.Command{
	Use: "image",
}

// TODO: Add flags for containerd, runc, cni, and kubernetes version
// TODO: Add image cleanup functionality

var imageBuildCmd = &cobra.Command{
	Use:   "build [name]",
	Short: "Build the packer image",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		name := args[0]

		command := []string{
			"image",
			"build",
			"--template",
			fmt.Sprintf("%s.pkr.hcl", name),
		}

		runner := container.New(cfg)
		err := runner.Run(ctx, command)
		cobra.CheckErr(err)
	},
}

var imageListCommand = &cobra.Command{
	Use:   "list",
	Short: "List simulator AMIs",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		ec2 := aws.EC2{}
		amis, err := ec2.List(ctx)
		cobra.CheckErr(err)

		table := tablewriter.NewWriter(os.Stdout)

		table.SetHeader([]string{
			"ID",
			"Name",
			"Created",
		})

		table.SetHeaderColor(
			tablewriter.Colors{tablewriter.Bold},
			tablewriter.Colors{tablewriter.Bold},
			tablewriter.Colors{tablewriter.Bold},
		)

		for _, ami := range amis {
			table.Append([]string{
				ami.ID,
				ami.Name,
				ami.CreationDate(),
			})
			table.SetRowLine(true)
		}
		table.Render()
	},
}

var imageDeleteCommand = &cobra.Command{
	Use:   "delete [ami id]",
	Short: "Delete a simulator AMI",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		id := args[0]
		ec2 := aws.EC2{}

		err := ec2.Delete(ctx, id)
		cobra.CheckErr(err)
	},
}

func init() {
	imageCmd.AddCommand(imageBuildCmd)
	imageCmd.AddCommand(imageListCommand)
	imageCmd.AddCommand(imageDeleteCommand)
	simulatorCmd.AddCommand(imageCmd)
}
