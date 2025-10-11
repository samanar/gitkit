package cmd

import (
	"fmt"
	"gitkit/git"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := git.CreateConfig()
		if err != nil {
			fmt.Printf("❌ Error creating config file ,%v\n", err)
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
