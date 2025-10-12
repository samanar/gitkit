/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/samanar/gitkit/git"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "git add",
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := git.NewGitCmdWithoutConfig()
		var files []string
		for i, arg := range args {
			if i > 0 {
				files = append(files, arg)
			}
		}
		gitCmd.Add(files...)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
