package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/samanar/gitkit/git"
	"github.com/spf13/cobra"
)

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Interactive branch selector and switcher",
	Long: `Display an interactive list of all branches with search functionality.
Select a branch with Enter to check it out.

Examples:
  gitkit branch           # Show interactive branch selector
  gitkit branch --all     # Include remote branches in selector`,
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := git.NewGitCmdWithoutConfig()

		// Get current branch to highlight it
		currentBranch := gitCmd.CurrentBranch()

		// Get list of branches
		branches := getAllBranches(gitCmd)

		if len(branches) == 0 {
			fmt.Println("❌ No branches found")
			return
		}

		// Create labels with current branch indicator
		branchLabels := make([]string, len(branches))
		selectedIdx := 0
		for i, branch := range branches {
			if branch == currentBranch {
				branchLabels[i] = fmt.Sprintf("✓ %s (current)", branch)
				selectedIdx = i
			} else {
				branchLabels[i] = branch
			}
		}

		// Create interactive selector with search
		searcher := func(input string, index int) bool {
			branch := branches[index]
			name := strings.Replace(strings.ToLower(branch), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)
			return strings.Contains(name, input)
		}

		prompt := promptui.Select{
			Label:             "Select a branch to checkout",
			Items:             branchLabels,
			CursorPos:         selectedIdx,
			Size:              15,
			Searcher:          searcher,
			StartInSearchMode: false,
		}

		idx, _, err := prompt.Run()
		if err != nil {
			// User cancelled (Ctrl+C or Esc)
			if err == promptui.ErrInterrupt {
				fmt.Println("\n👋 Branch selection cancelled")
				return
			}
			fmt.Fprintf(os.Stderr, "❌ Selection failed: %v\n", err)
			return
		}

		selectedBranch := branches[idx]

		// Don't checkout if already on this branch
		if selectedBranch == currentBranch {
			fmt.Printf("ℹ️  Already on branch '%s'\n", selectedBranch)
			return
		}

		// Check if it's a remote branch that needs to be tracked
		if strings.HasPrefix(selectedBranch, "remotes/") {
			// Extract branch name from remotes/origin/branch-name
			parts := strings.Split(selectedBranch, "/")
			if len(parts) >= 3 {
				localBranch := strings.Join(parts[2:], "/")

				// Check if local branch already exists
				if gitCmd.BranchExists(localBranch) {
					fmt.Printf("🔄 Switching to local branch '%s'\n", localBranch)
					gitCmd.Checkout(localBranch)
					fmt.Printf("✅ Switched to branch '%s'\n", localBranch)
				} else {
					// Create and track the remote branch
					fmt.Printf("🔄 Creating local branch '%s' tracking '%s'\n", localBranch, selectedBranch)
					output, err := gitCmd.Run("checkout", "-b", localBranch, "--track", selectedBranch)
					if err != nil {
						fmt.Fprintf(os.Stderr, "❌ Failed to checkout remote branch: %v\n%s\n", err, output)
						return
					}
					fmt.Printf("✅ Created and switched to branch '%s'\n", localBranch)
				}
			}
			return
		}

		// Checkout the selected branch
		fmt.Printf("🔄 Switching to branch '%s'\n", selectedBranch)
		output, err := gitCmd.Run("checkout", selectedBranch)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to checkout branch: %v\n%s\n", err, output)
			return
		}

		fmt.Printf("✅ Switched to branch '%s'\n", selectedBranch)
	},
}

// getAllBranches gets both local and remote branches
func getAllBranches(gitCmd git.GitCmd) []string {
	output := gitCmd.RunMust("branch", "--all")
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var branches []string

	for _, line := range lines {
		branch := strings.TrimSpace(line)
		branch = strings.TrimPrefix(branch, "* ") // remove asterisk from current branch

		// Skip HEAD references
		if strings.Contains(branch, "HEAD ->") {
			continue
		}

		if branch != "" {
			branches = append(branches, branch)
		}
	}

	return branches
}

func init() {
	rootCmd.AddCommand(branchCmd)
}
