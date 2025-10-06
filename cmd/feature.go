package cmd

import (
	"fmt"
	"gitkit/config"
	"gitkit/git"
	"os"

	"github.com/spf13/cobra"
)

// featureStartCmd represents the 'feature start' command
var featureStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new feature",
	Long:  `Start a new feature branch or process in gitkit.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not load .gitkit.yml: %v\n", err)
			os.Exit(1)
		}

		featurePrefix := cfg.Prefixes.Feature
		developBranch := cfg.Branches.Develop
		featureName := args[0]

		git.SyncRemoteBranch(developBranch)

		// Use reusable method to create the branch
		if err := git.CreatePrefixedBranch(developBranch, featurePrefix, featureName); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Feature branch '%s%s' started from '%s'.\n", featurePrefix, featureName, developBranch)
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
