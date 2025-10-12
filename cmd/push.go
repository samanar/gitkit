package cmd

import (
	"github.com/samanar/gitkit/git"

	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "git push",
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := git.NewGitCmdWithoutConfig()
		gitCmd.Push()
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

}
