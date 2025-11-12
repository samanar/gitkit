package cmd

import (
	"fmt"
	"os"

	"github.com/samanar/gitkit/git"
	"github.com/spf13/cobra"
)

var (
	cloneDepth   int
	cloneBranch  string
	cloneRecurse bool
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone <repository> [directory]",
	Short: "git clone with optional gitkit initialization",
	Long: `Clone a Git repository and optionally initialize gitkit configuration.
    
Examples:
  gitkit clone https://github.com/user/repo
  gitkit clone https://github.com/user/repo my-project
  gitkit clone --branch develop https://github.com/user/repo
  gitkit clone --depth 1 https://github.com/user/repo`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		repository := args[0]
		var directory string
		if len(args) == 2 {
			directory = args[1]
		}

		gitCmd := git.NewGitCmdWithoutConfig()

		// Build clone arguments
		var cloneArgs []string
		cloneArgs = append(cloneArgs, "clone")

		if cloneDepth > 0 {
			cloneArgs = append(cloneArgs, "--depth", fmt.Sprintf("%d", cloneDepth))
		}

		if cloneBranch != "" {
			cloneArgs = append(cloneArgs, "--branch", cloneBranch)
		}

		if cloneRecurse {
			cloneArgs = append(cloneArgs, "--recurse-submodules")
		}

		cloneArgs = append(cloneArgs, repository)

		if directory != "" {
			cloneArgs = append(cloneArgs, directory)
		}

		fmt.Printf("🔄 Cloning repository: %s\n", repository)
		output := gitCmd.RunMust(cloneArgs...)
		if output != "" {
			fmt.Print(output)
		}

		fmt.Println("✅ Repository cloned successfully!")

		// Determine target directory
		targetDir := directory
		if targetDir == "" {
			// Extract directory name from repository URL
			targetDir = extractRepoName(repository)
		}

		// Change to cloned directory
		if err := os.Chdir(targetDir); err != nil {
			return
		}

		// Check if .gitkit.yaml already exists
		if _, err := os.Stat(".gitkit.yaml"); err == nil {
			return
		}

		// Ask if user wants to initialize gitkit
		fmt.Print("\nDo you want to initialize gitkit in this repository? (y/N): ")
		var response string
		fmt.Scanln(&response)

		if response == "y" || response == "Y" || response == "yes" {
			// Initialize gitkit
			git.NewGitCommandWithForceRewriteConfig()
			fmt.Println("✅ GitKit initialized!")
		}
	},
}

// extractRepoName extracts repository name from URL
func extractRepoName(repoURL string) string {
	// Remove .git suffix if present
	name := repoURL
	if len(name) > 4 && name[len(name)-4:] == ".git" {
		name = name[:len(name)-4]
	}

	// Get last part after /
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '/' {
			return name[i+1:]
		}
	}
	return name
}

func init() {
	cloneCmd.Flags().IntVar(&cloneDepth, "depth", 0, "Create a shallow clone with a history truncated to the specified number of commits")
	cloneCmd.Flags().StringVarP(&cloneBranch, "branch", "b", "", "Clone a specific branch instead of the default")
	cloneCmd.Flags().BoolVar(&cloneRecurse, "recurse-submodules", false, "Initialize and clone submodules within the cloned repository")

	rootCmd.AddCommand(cloneCmd)
}
