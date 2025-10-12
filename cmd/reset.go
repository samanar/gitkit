package cmd

import (
	"fmt"
	"os"

	"github.com/samanar/gitkit/git"

	"github.com/spf13/cobra"
)

var (
	resetHard  bool
	resetSoft  bool
	resetMixed bool
	resetKeep  bool
	resetMerge bool
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "git reset",
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := git.NewGitCmdWithoutConfig()

		var modes []string
		if resetHard {
			modes = append(modes, "--hard")
		}
		if resetSoft {
			modes = append(modes, "--soft")
		}
		if resetMixed {
			modes = append(modes, "--mixed")
		}
		if resetKeep {
			modes = append(modes, "--keep")
		}
		if resetMerge {
			modes = append(modes, "--merge")
		}

		if len(modes) > 1 {
			fmt.Fprintln(os.Stderr, "Please specify only one reset mode flag.")
			os.Exit(1)
		}

		var resetArgs []string
		if len(modes) == 1 {
			resetArgs = append(resetArgs, modes[0])
		}
		resetArgs = append(resetArgs, args...)

		gitCmd.Reset(resetArgs...)
	},
}

func init() {
	resetCmd.Flags().BoolVar(&resetHard, "hard", false, "Reset HEAD, index, and working tree, discarding changes")
	resetCmd.Flags().BoolVar(&resetSoft, "soft", false, "Reset HEAD only, keeping index and working tree")
	resetCmd.Flags().BoolVar(&resetMixed, "mixed", false, "Reset HEAD and index, keeping working tree (default mode)")
	resetCmd.Flags().BoolVar(&resetKeep, "keep", false, "Reset HEAD and index like --hard but keep local changes that don't conflict")
	resetCmd.Flags().BoolVar(&resetMerge, "merge", false, "Reset HEAD, index, and working tree like --hard but keep unmerged entries")
	rootCmd.AddCommand(resetCmd)
}
