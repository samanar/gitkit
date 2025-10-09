package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stderr.String(), fmt.Errorf("git %v failed: %v", args, err)
	}

	return out.String(), nil
}

func RunMust(args ...string) string {
	output, err := Run(args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ git %v failed:\n%s\n", args, output)
		os.Exit(1)
	}
	return output
}

func Push() {
	RunMust("push")
}

func Pull() {
	RunMust("pull", "--rebase")
}

func Fetch() {
	RunMust("fetch")
}

func IsClean() bool {
	output := RunMust("status", "--porcelain")
	return strings.TrimSpace(output) == ""
}

func CreateBranch(branch string) {
	RunMust("checkout", "-b", branch)
}

func CreateAndPushBranch(branch string) {
	CreateBranch(branch)
	PushWithSetUpstream(branch)
}

func PushWithSetUpstream(branch string) {
	RunMust("push", "--set-upstream", "origin", branch)
}

func Merge(branch string) {
	RunMust("merge", "--no-ff", branch)
}

func MergeWithCommitMessage(branch, commit string) {
	RunMust("merge", "--no-ff", branch, "-m", commit)
}

func MergeSquash(branch string) {
	RunMust("merge", "--squash", branch)
}

func CommitAll(message string) {
	RunMust("add", "-A")
	RunMust("commit", "-m", message)
}

func Tag(tag, message string) {
	RunMust("tag", "-a", tag, "-m", message)
}

func RootDir() (string, error) {
	out, err := Run("rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func SyncRemoteBranch(branch string) {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Could not load .gitkit.yml: %v\n", err)
		os.Exit(1)
	}
	if _, err := Run("pull", "--rebase", cfg.Remote, branch); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to pull latest '%s': %v\n", branch, err)
		os.Exit(1)
	}
}

// CreatePrefixedBranch creates a new branch from baseBranch with prefix and branchName.
func CreatePrefixedBranch(baseBranch, prefix, branchName string) error {
	// Checkout base branch
	if _, err := Run("checkout", baseBranch); err != nil {
		return fmt.Errorf("failed to checkout '%s': %w", baseBranch, err)
	}
	// Create new branch from base
	newBranch := prefix + branchName
	if _, err := Run("checkout", "-b", newBranch, baseBranch); err != nil {
		return fmt.Errorf("failed to create branch '%s': %w", newBranch, err)
	}
	return nil
}

func MergeBranchToBase(baseBranch, branch string) {
	if !BranchExists(branch) {
		fmt.Fprintf(os.Stderr, "❌ branch '%s' does not exist.\n", branch)
		os.Exit(1)
	}
	Checkout(baseBranch)
	Pull()
	Merge(branch)
	DeleteBranchSafe(branch)
	Push()
}
