package cmd

import (
	"gitkit/git"

	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "git pull",
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := git.NewGitCmdWithoutConfig()
		gitCmd.Pull()
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
