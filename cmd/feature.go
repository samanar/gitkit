package cmd

import (
	"gitkit/git"

	"github.com/spf13/cobra"
)

var featureCmd = &cobra.Command{
	Use:     "feature",
	Aliases: []string{"f", "feat"},
	Short:   "Feature related commands",
	Long:    `Commands to manage features in gitkit.`,
}

// featureStartCmd represents the 'feature start' command
var featureStartCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start a new feature",
	Aliases: []string{"s", "new", "begin"},
	Long:    `Start a new feature branch or process in gitkit.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		gitCmd := git.NewGitCommandWithConfig(false)
		gitCmd.StartBranch("feature", branchName)
	},
}

// featureEndCmd represents the 'feature end' command
var featureEndCmd = &cobra.Command{
	Use:     "finish",
	Short:   "Finish the current feature",
	Aliases: []string{"f", "end", "complete"},
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := git.NewGitCommandWithConfig(false)
		var branchName string = gitCmd.CurrentBranch()
		if len(args) == 1 {
			branchName = args[0]
		}
		gitCmd.FinishBranch("feature", branchName)
	},
}

func init() {
	featureCmd.AddCommand(featureStartCmd)
	featureCmd.AddCommand(featureEndCmd)
	rootCmd.AddCommand(featureCmd)
}
