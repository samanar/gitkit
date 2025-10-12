package cmd

import (
	"github.com/samanar/gitkit/git"

	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "git reset",
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := git.NewGitCmdWithoutConfig()
		gitCmd.Reset(args...)
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
