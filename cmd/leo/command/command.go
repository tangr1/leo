package command

import (
	"leo/cmd/leo/config"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	rootCmd := &cobra.Command{Use: "leo"}
	rootCmd.AddCommand(VersionCommand())
	config.InitFlag(rootCmd)
	return rootCmd
}
