/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gitkit/git"

	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "git commit -m",
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := git.NewGitCmdWithoutConfig()
		if len(args) == 0 {
			fmt.Println("Please provide a commit message.")
			return
		}
		message := ""
		for i, arg := range args {
			if i > 0 {
				message += " "
			}
			message += arg
		}
		gitCmd.CommitAll(message)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
