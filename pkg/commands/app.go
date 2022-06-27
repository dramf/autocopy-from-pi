package commands

import (
	"github.com/spf13/cobra"
)

func NewApp(version string) *cobra.Command {
	rootCmd.Version = version
	return rootCmd
}

var rootCmd = &cobra.Command{
	Use:   "autocopy",
	Short: "a simple CLI tool for auto copying data from mounted devices",
	Run:   func(cmd *cobra.Command, args []string) {},
}
