/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/samanar/gitkit/git"

	"github.com/spf13/cobra"
)

// finishCmd represents the finish command
var finishCmd = &cobra.Command{
	Use:   "finish",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := git.NewGitCommandWithConfig(false)
		currentBranch := gitCmd.CurrentBranch()
		branchType := ""
		for key, branch := range gitCmd.Config.Prefixes {
			if strings.HasPrefix(currentBranch, branch.Name) {
				branchType = key
				break
			}
		}
		if branchType != "" {
			gitCmd.FinishBranch(branchType, currentBranch)
		}
	},
}

func init() {
	rootCmd.AddCommand(finishCmd)
}
