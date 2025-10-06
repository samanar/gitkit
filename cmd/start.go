package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// startCmd represents the 'start' command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the gitkit feature",
	Long:  `Start the gitkit feature. This command initializes or launches the gitkit process.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gitkit feature started!")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
