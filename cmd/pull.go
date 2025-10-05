package cmd

import (
	"fmt"
	"gitkit/git"

	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📥 Pulling latest changes...")
		git.Pull()
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

}
