package cmd

import (
	"fmt"
	"gitkit/git"

	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📤 Pushing local commits...")
		git.Push()
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

}
