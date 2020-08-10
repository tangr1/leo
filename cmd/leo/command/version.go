package command

import (
	"github.com/fatih/color"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/cobra"
)

func VersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Binary version",
		Long:  "Binary version",
		Args:  cobra.MinimumNArgs(0),
		Run:   version,
	}
}

func version(cmd *cobra.Command, args []string) {
	box := packr.New("version_box", "version")
	version, err := box.FindString("version.txt")
	if err != nil {
		color.Red("Failed to get version: %v", err)
	}
	color.Green(version)
}
