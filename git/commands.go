package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (g *GitCmd) Run(args ...string) (string, error) {
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

func (g *GitCmd) RunMust(args ...string) string {
	output, err := g.Run(args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ git %v failed:\n%s\n", args, output)
		os.Exit(1)
	}
	return output
}

func (g *GitCmd) Push() {
	fmt.Println("📤 Pushing local commits...")
	if g.Config.Remote == "" {
		g.RunMust("push")
		return
	}
	g.RunMust("push", "-u", g.Config.Remote, g.CurrentBranch())
}

func (g *GitCmd) Pull() {
	fmt.Println("📥 Pulling latest changes...")
	g.RunMust("pull", "--rebase")
}

func (g *GitCmd) Sync(branch string) {
	if output, err := g.Run("pull", "--rebase", g.Config.Remote, branch); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to pull latest '%s': %v\n", branch, err)
		fmt.Println(output)
		os.Exit(1)
	}
}

func (g *GitCmd) Fetch() {
	g.RunMust("fetch")
}

func (g *GitCmd) IsClean() bool {
	output := g.RunMust("status", "--porcelain")
	return strings.TrimSpace(output) == ""
}

func (g *GitCmd) CreateBranch(branch string) {
	g.RunMust("checkout", "-b", branch)
}

func (g *GitCmd) PushWithSetUpstream(branch string) {
	g.RunMust("push", "--set-upstream", "origin", branch)
}

func (g *GitCmd) Merge(branch string) error {
	_, err := g.Run("merge", "--no-ff", branch)
	return err
}

func (g *GitCmd) MergeWithCommitMessage(branch, commit string) {
	g.RunMust("merge", "--no-ff", branch, "-m", commit)
}

func (g *GitCmd) Status() {
	output := g.RunMust("status", "--porcelain")
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

func (g *GitCmd) Add(files ...string) {
	if len(files) == 0 {
		g.RunMust("add", "-A")
		return
	}
	args := append([]string{"add"}, files...)
	g.RunMust(args...)
}

func (g *GitCmd) CommitAll(message string) {
	g.RunMust("add", "-A")
	g.RunMust("commit", "-m", message)
}

func (g *GitCmd) Tag(tag, message string) {
	g.RunMust("tag", "-a", tag, "-m", message)
}

func (g *GitCmd) RootDir() (string, error) {
	out, err := g.Run("rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// CreatePrefixedBranch creates a new branch from baseBranch with prefix and branchName.
func (g *GitCmd) CreatePrefixedBranch(baseBranch, prefix, branchName string) error {
	// Checkout base branch
	if _, err := g.Run("checkout", baseBranch); err != nil {
		return fmt.Errorf("failed to checkout '%s': %w", baseBranch, err)
	}
	// Create new branch from base
	newBranch := prefix + branchName
	if _, err := g.Run("checkout", "-b", newBranch, baseBranch); err != nil {
		return fmt.Errorf("failed to create branch '%s': %w", newBranch, err)
	}
	return nil
}

func (g *GitCmd) MergeBranchToBase(baseBranch, branch string) error {
	if !g.BranchExists(branch) {
		fmt.Fprintf(os.Stderr, "❌ branch '%s' does not exist.\n", branch)
		os.Exit(1)
	}
	g.Checkout(baseBranch)
	g.Pull()
	err := g.Merge(branch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to merge branch '%s': %v\n", branch, err)
		return err
	}
	g.DeleteBranchSafe(branch)
	g.Push()
	return nil
}

func (g *GitCmd) Reset(args ...string) {
	resetArgs := append([]string{"reset"}, args...)
	output := g.RunMust(resetArgs...)
	if strings.TrimSpace(output) != "" {
		fmt.Print(output)
	}
}

func (g *GitCmd) Restore(args ...string) {
	restoreArgs := append([]string{"restore"}, args...)
	output := g.RunMust(restoreArgs...)
	if strings.TrimSpace(output) != "" {
		fmt.Print(output)
	}
}

func (g *GitCmd) Stash(args ...string) {
	stashArgs := append([]string{"stash"}, args...)
	output := g.RunMust(stashArgs...)
	if strings.TrimSpace(output) != "" {
		fmt.Print(output)
	}
}

func (g *GitCmd) StashPop(args ...string) {
	stashArgs := append([]string{"stash", "pop"}, args...)
	output := g.RunMust(stashArgs...)
	if strings.TrimSpace(output) != "" {
		fmt.Print(output)
	}
}

// Clone clones a repository with the given arguments
func (g *GitCmd) Clone(args ...string) {
	cloneArgs := append([]string{"clone"}, args...)
	output := g.RunMust(cloneArgs...)
	if strings.TrimSpace(output) != "" {
		fmt.Print(output)
	}
}
