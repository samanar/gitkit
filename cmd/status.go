/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gitkit/git"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "git status",

	Run: func(cmd *cobra.Command, args []string) {
		git.Status()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
