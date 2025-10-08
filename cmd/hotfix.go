package cmd

import (
	"gitkit/git"

	"github.com/spf13/cobra"
)

var hotFixCmd = &cobra.Command{
	Use:     "hotfix",
	Aliases: []string{"hot"},
	Short:   "hotFix related commands",
	Long:    `Commands to manage hotFixes in gitkit.`,
}

var hotFixStartCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start a new feature",
	Aliases: []string{"s", "new", "begin"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		git.StartBranch("hotFix", branchName)
	},
}

var hotFixEndCmd = &cobra.Command{
	Use:     "finish",
	Short:   "Finish the current hotFix",
	Aliases: []string{"f", "end", "complete"},
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var branchName string = ""
		if len(args) == 1 {
			branchName = args[0]
		}
		git.FinishBranch("hotFix", branchName)
	},
}

func init() {
	hotFixCmd.AddCommand(hotFixStartCmd)
	hotFixCmd.AddCommand(hotFixEndCmd)
	rootCmd.AddCommand(hotFixCmd)
}
