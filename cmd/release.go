package cmd

import (
	"fmt"
	"gitkit/git"
	"os"

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
		var branchName string = git.CurrentBranch()
		if len(args) == 1 {
			branchName = args[0]
		}

		cfg, err := git.LoadConfig()

		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not load .gitkit.yml: %v\n", err)
			os.Exit(1)
		}
		prefixCfg, ok := cfg.Prefixes["release"]
		if !ok {
			fmt.Printf("❌ Unknown type: %s\n", "release")
			return
		}
		branchName = git.RemovePrefix(branchName, prefixCfg.Name)
		branchName = prefixCfg.Name + branchName
		if !git.BranchExists(branchName) {
			fmt.Fprintf(os.Stderr, "❌ branch '%s' does not exist.\n", branchName)
			os.Exit(1)
		}
		developBranch := cfg.Branches.Develop
		mainBranch := cfg.Branches.Main

		// # 1. Ensure you're up to date
		git.Checkout(developBranch)
		git.Pull()

		// # 2. Checkout the release branch
		git.Checkout(branchName)

		// # 3. Merge the release into main (production)
		git.Checkout(mainBranch)
		git.Pull()
		git.MergeWithCommitMessage(branchName, fmt.Sprintf("Merge release %s into %s", branchName, mainBranch))

		// # 4. Tag the release
		git.Tag(git.RemovePrefix(branchName, prefixCfg.Name), branchName)

		// # 5. Merge release back into develop
		git.Checkout(developBranch)
		git.MergeWithCommitMessage(branchName, fmt.Sprintf("Merge release %s into %s", branchName, mainBranch))

		// # 6. Delete the release branch
		git.DeleteBranchSafe(branchName)

		// # 7. Push everything (branches + tags)
		// git push origin develop main --tags
		git.RunMust("push", "origin", developBranch, mainBranch, "--tags")
		fmt.Printf("✅ Finished release: %s'.\n", branchName)
	},
}

func init() {
	releaseCmd.AddCommand(releaseStartCmd)
	releaseCmd.AddCommand(releaseEndCmd)
	rootCmd.AddCommand(releaseCmd)
}
