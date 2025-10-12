package cmd

import (
	"fmt"

	"github.com/samanar/gitkit/git"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔄 Syncing with remote...")
		gitCmd := git.NewGitCmdWithoutConfig()
		gitCmd.Pull()
		gitCmd.Push()
		fmt.Println("✅ Sync complete.")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
