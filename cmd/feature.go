package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// featureStartCmd represents the 'feature start' command
var featureStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new feature",
	Long:  `Start a new feature branch or process in gitkit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Feature started!")
	},
}

// featureEndCmd represents the 'feature end' command
var featureEndCmd = &cobra.Command{
	Use:   "end",
	Short: "End the current feature",
	Long:  `End the current feature branch or process in gitkit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Feature ended!")
	},
}

var featureCmd = &cobra.Command{
	Use:   "feature",
	Short: "Feature related commands",
	Long:  `Commands to manage features in gitkit.`,
}

func init() {
	featureCmd.AddCommand(featureStartCmd)
	featureCmd.AddCommand(featureEndCmd)
	rootCmd.AddCommand(featureCmd)
}
