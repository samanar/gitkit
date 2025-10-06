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

func MergeSquash(branch string) {
	RunMust("merge", "--squash", branch)
}

func CommitAll(message string) {
	RunMust("add", "-A")
	RunMust("commit", "-m", message)
}

func Tag(tag string) {
	RunMust("tag", tag)
	RunMust("push", "origin", tag)
}
