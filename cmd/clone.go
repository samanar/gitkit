package cmd

import (
	"fmt"
	"gitkit/git"

	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone [repo URL]",
	Short: "Clone a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := args[0]
		fmt.Println("📥 Cloning repository:", repoURL)
		if err := git.Clone(repoURL); err != nil {
			fmt.Println("❌ Clone failed:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
