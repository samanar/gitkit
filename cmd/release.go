package cmd

import (
	"gitkit/git"

	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:     "release",
	Aliases: []string{"version"},
	Short:   "Release related commands",
}

var releaseStartCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"s", "new", "begin"},
	Short:   "Start a new release",
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		git.StartBranch("feature", branchName)
	},
}

// releaseEndCmd represents the 'release end' command
var releaseEndCmd = &cobra.Command{
	Use:     "end",
	Aliases: []string{"f", "end", "complete"},
	Short:   "End the current release",
	Run: func(cmd *cobra.Command, args []string) {
		var branchName string = ""
		if len(args) == 1 {
			branchName = args[0]
		}
		git.FinishBranch("feature", branchName)
	},
}

func init() {
	releaseCmd.AddCommand(releaseStartCmd)
	releaseCmd.AddCommand(releaseEndCmd)
	rootCmd.AddCommand(releaseCmd)
}
