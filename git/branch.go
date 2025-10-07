package git

import "strings"

func CurrentBranch() string {
	output := RunMust("branch", "--show-current")
	return output[:len(output)-1]
}

func Branches() []string {
	output := RunMust("branch", "--list")
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var branches []string
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		branch = strings.TrimPrefix(branch, "* ") // remove asterisk from current branch
		if branch != "" {
			branches = append(branches, branch)
		}
	}
	return branches
}

func BranchesWithPrefix(prefix string) []string {
	branches := Branches()
	var matches []string

	for _, branch := range branches {
		if strings.HasPrefix(branch, prefix) {
			matches = append(matches, branch)
		}
	}
	return matches
}

func DeleteBranch(branch string, force bool) {
	flag := "-d"
	if force {
		flag = "-D"
	}
	RunMust("branch", flag, branch)
}

func DeleteBranchSafe(branch string) {
	DeleteBranch(branch, false)
}

func Checkout(branch string) {
	RunMust("checkout", branch)
	Pull()
}

func RemovePrefix(branch, prefix string) string {
	if strings.HasPrefix(branch, prefix) {
		return branch[len(prefix):]
	}
	return branch
}

func BranchExists(branch string) bool {
	branches := Branches()
	for _, b := range branches {
		if b == branch {
			return true
		}
	}
	return false
}
