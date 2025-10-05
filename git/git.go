package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

func RunGitCommand(args ...string) (string, error) {
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
