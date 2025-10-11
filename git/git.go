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
	cfg, err := LoadConfig()
	if err != nil {
		RunMust("push")
	} else {
		RunMust("push", "-u", cfg.Remote, CurrentBranch())
	}
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

func Merge(branch string) error {
	_, err := Run("merge", "--no-ff", branch)
	return err
}

func MergeWithCommitMessage(branch, commit string) {
	RunMust("merge", "--no-ff", branch, "-m", commit)
}

func MergeSquash(branch string) {
	RunMust("merge", "--squash", branch)
}

func Status() {
	output := RunMust("status", "--porcelain")
	if output == "" {
		fmt.Println("✅ Working directory clean!")
		return
	}
	fmt.Println("🔍 Git Status:")
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		status := line[:2]
		file := strings.TrimSpace(line[2:])
		var symbol, desc string
		switch status {
		case " M":
			symbol, desc = "✏️", "Modified"
		case "A ":
			symbol, desc = "➕", "Added"
		case " D":
			symbol, desc = "❌", "Deleted"
		case "??":
			symbol, desc = "🆕", "Untracked"
		case "R ":
			symbol, desc = "🔀", "Renamed"
		case "C ":
			symbol, desc = "📋", "Copied"
		default:
			symbol, desc = "❓", "Other"
		}
		fmt.Printf("%s %-10s %s\n", symbol, desc, file)
	}
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

func MergeBranchToBase(baseBranch, branch string) error {
	if !BranchExists(branch) {
		fmt.Fprintf(os.Stderr, "❌ branch '%s' does not exist.\n", branch)
		os.Exit(1)
	}
	Checkout(baseBranch)
	Pull()
	err := Merge(branch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to merge branch '%s': %v\n", branch, err)
		return err
	}
	DeleteBranchSafe(branch)
	Push()
	return nil
}
