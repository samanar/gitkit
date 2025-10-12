package cmd

import (
	"github.com/samanar/gitkit/git"

	"github.com/spf13/cobra"
)

var stashCmd = &cobra.Command{
	Use:   "stash",
	Short: "git stash",
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := git.NewGitCmdWithoutConfig()
		gitCmd.Stash(args...)
	},
}

var stashPopCmd = &cobra.Command{
	Use:   "pop",
	Short: "git stash pop",
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := git.NewGitCmdWithoutConfig()
		gitCmd.StashPop(args...)
	},
}

func init() {
	stashCmd.AddCommand(stashPopCmd)
	rootCmd.AddCommand(stashCmd)
}
