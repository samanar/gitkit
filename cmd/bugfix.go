package cmd

import (
	"gitkit/git"

	"github.com/spf13/cobra"
)

var bugFixCmd = &cobra.Command{
	Use:     "bugfix",
	Aliases: []string{"fix", "bug"},
	Short:   "Bugfix related commands",
}

var bugFixStartCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start a new bugFix",
	Aliases: []string{"s", "new", "begin"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		git.StartBranch("bugFix", branchName)
	},
}

var bugFixFinishCmd = &cobra.Command{
	Use:     "finish",
	Aliases: []string{"f", "end", "complete"},
	Short:   "Finish the current Bugfix",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var branchName string = ""
		if len(args) == 1 {
			branchName = args[0]
		}
		git.FinishBranch("bugFix", branchName)
	},
}

func init() {
	bugFixCmd.AddCommand(bugFixStartCmd)
	bugFixCmd.AddCommand(bugFixFinishCmd)
	rootCmd.AddCommand(bugFixCmd)
}
