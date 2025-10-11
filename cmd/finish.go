/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gitkit/git"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// finishCmd represents the finish command
var finishCmd = &cobra.Command{
	Use:   "finish",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("finish called")
		currentBranch := git.CurrentBranch()
		cfg, err := git.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not load .gitkit.yml: %v\n", err)
			os.Exit(1)
		}
		branchType := ""
		for key, branch := range cfg.Prefixes {
			if strings.HasPrefix(currentBranch, branch.Name) {
				branchType = key
				break
			}
		}
		if branchType != "" {
			git.FinishBranch(branchType, currentBranch)

		}

	},
}

func init() {
	rootCmd.AddCommand(finishCmd)
}
