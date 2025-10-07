package cmd

import (
	"fmt"
	"gitkit/git"
	"os"

	"github.com/spf13/cobra"
)

var bugFixCmd = &cobra.Command{
	Use:   "bugfix",
	Short: "Bugfix related commands",
}

var bugFixStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new bugFix",
	Long:  `Start a new feature branch or process in gitkit.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := git.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not load .gitkit.yml: %v\n", err)
			os.Exit(1)
		}

		bugFixPrefix := cfg.Prefixes.BugFix
		developBranch := cfg.Branches.Develop
		bugFixName := args[0]
		bugFixName = git.RemovePrefix(bugFixName, bugFixPrefix)

		git.SyncRemoteBranch(developBranch)

		// Use reusable method to create the branch
		if err := git.CreatePrefixedBranch(developBranch, bugFixPrefix, bugFixName); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Feature branch '%s%s' started from '%s'.\n", bugFixPrefix, bugFixName, developBranch)
	},
}

var bugFixFinishCmd = &cobra.Command{
	Use:   "finish",
	Short: "Finish the current feature",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := git.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not load .gitkit.yml: %v\n", err)
			os.Exit(1)
		}
		bugFixPrefix := cfg.Prefixes.BugFix
		developBranch := cfg.Branches.Develop
		var branch string
		if len(args) == 1 {
			branch = args[0]
		} else {
			branch = git.CurrentBranch()
		}
		branch = git.RemovePrefix(branch, bugFixPrefix)
		branch = bugFixPrefix + branch
		git.MergeBranchToBase(developBranch, branch)

		fmt.Printf("✅ BugFix branch '%s' finished and merged into '%s'.\n", branch, developBranch)
	},
}

func init() {
	bugFixCmd.AddCommand(bugFixStartCmd)
	bugFixCmd.AddCommand(bugFixFinishCmd)
	rootCmd.AddCommand(bugFixCmd)
}
