package cmd

import (
	"github.com/samanar/gitkit/git"

	"github.com/spf13/cobra"
)

var (
	restoreStaged   bool
	restoreWorktree bool
	restoreSource   string
	restorePaths    []string
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "git restore",
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := git.NewGitCmdWithoutConfig()

		var restoreArgs []string
		if restoreStaged {
			restoreArgs = append(restoreArgs, "--staged")
		}
		if restoreWorktree {
			restoreArgs = append(restoreArgs, "--worktree")
		}
		if restoreSource != "" {
			restoreArgs = append(restoreArgs, "--source", restoreSource)
		}
		restoreArgs = append(restoreArgs, restorePaths...)
		restoreArgs = append(restoreArgs, args...)

		gitCmd.Restore(restoreArgs...)
	},
}

func init() {
	restoreCmd.Flags().BoolVar(&restoreStaged, "staged", false, "Restore paths in the index")
	restoreCmd.Flags().BoolVar(&restoreWorktree, "worktree", false, "Restore paths in the working tree (default behaviour)")
	restoreCmd.Flags().StringVar(&restoreSource, "source", "", "Restore paths from the given tree-ish")
	restoreCmd.Flags().StringSliceVar(&restorePaths, "file", nil, "Specific file(s) to restore")
	rootCmd.AddCommand(restoreCmd)
}
