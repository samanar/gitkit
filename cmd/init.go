package cmd

import (
	"gitkit/git"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "creating/overwriting config file for gitkit",
	RunE: func(cmd *cobra.Command, args []string) error {
		git.NewGitCommandWithConfig(true)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
