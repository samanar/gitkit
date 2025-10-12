package cmd

import (
	"github.com/samanar/gitkit/git"

	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "git restore",
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := git.NewGitCmdWithoutConfig()
		gitCmd.Restore(args...)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
