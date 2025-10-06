package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// releaseStartCmd represents the 'release start' command
var releaseStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new release",
	Long:  `Start a new release branch or process in gitkit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Release started!")
	},
}

// releaseEndCmd represents the 'release end' command
var releaseEndCmd = &cobra.Command{
	Use:   "end",
	Short: "End the current release",
	Long:  `End the current release branch or process in gitkit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Release ended!")
	},
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Release related commands",
	Long:  `Commands to manage releases in gitkit.`,
}

func init() {
	releaseCmd.AddCommand(releaseStartCmd)
	releaseCmd.AddCommand(releaseEndCmd)
	rootCmd.AddCommand(releaseCmd)
}
