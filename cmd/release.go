package cmd

import (
	"fmt"
	"os"

	"github.com/samanar/gitkit/git"

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
		gitCmd := git.NewGitCommandWithConfig(false)
		branchName := args[0]
		gitCmd.StartBranch("release", branchName)
	},
}

// releaseEndCmd represents the 'release end' command
var releaseEndCmd = &cobra.Command{
	Use:     "end",
	Aliases: []string{"f", "end", "complete"},
	Short:   "End the current release",
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := git.NewGitCommandWithConfig(false)
		var branchName string = gitCmd.CurrentBranch()
		if len(args) == 1 {
			branchName = args[0]
		}

		prefixCfg, ok := gitCmd.Config.Prefixes["release"]
		if !ok {
			fmt.Printf("❌ Unknown type: %s\n", "release")
			return
		}
		branchName = gitCmd.RemovePrefix(branchName, prefixCfg.Name)
		branchName = prefixCfg.Name + branchName
		if !gitCmd.BranchExists(branchName) {
			fmt.Fprintf(os.Stderr, "❌ branch '%s' does not exist.\n", branchName)
			os.Exit(1)
		}
		developBranch := gitCmd.Config.Branches.Develop
		mainBranch := gitCmd.Config.Branches.Main

		// # 1. Ensure you're up to date
		gitCmd.Checkout(developBranch)
		gitCmd.Pull()

		// # 2. Checkout the release branch
		gitCmd.Checkout(branchName)

		// # 3. Merge the release into main (production)
		gitCmd.Checkout(mainBranch)
		gitCmd.Pull()
		gitCmd.MergeWithCommitMessage(branchName, fmt.Sprintf("Merge release %s into %s", branchName, mainBranch))

		// # 4. Tag the release
		gitCmd.Tag(gitCmd.RemovePrefix(branchName, prefixCfg.Name), branchName)

		// # 5. Merge release back into develop
		gitCmd.Checkout(developBranch)
		gitCmd.MergeWithCommitMessage(branchName, fmt.Sprintf("Merge release %s into %s", branchName, mainBranch))

		// # 6. Delete the release branch
		gitCmd.DeleteBranchSafe(branchName)

		// # 7. Push everything (tbranches + tags)
		// gitCmd push origin develop main --tags
		gitCmd.RunMust("push", "origin", developBranch, mainBranch, "--tags")
		fmt.Printf("✅ Finished release: %s'.\n", branchName)
	},
}

func init() {
	releaseCmd.AddCommand(releaseStartCmd)
	releaseCmd.AddCommand(releaseEndCmd)
	rootCmd.AddCommand(releaseCmd)
}
