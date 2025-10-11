package git

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

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

func StartBranch(branchType, branchName string) {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Could not load .gitkit.yml: %v\n", err)
		os.Exit(1)
	}
	prefixCfg, ok := cfg.Prefixes[branchType]
	if !ok {
		fmt.Printf("❌ Unknown type: %s\n", branchType)
		return
	}
	branchName = RemovePrefix(branchName, prefixCfg.Name)

	SyncRemoteBranch(prefixCfg.Base)

	// // Use reusable method to create the branch
	if err := CreatePrefixedBranch(prefixCfg.Base, prefixCfg.Name, branchName); err != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Feature branch '%s%s' started from '%s'.\n", prefixCfg.Name, branchName, prefixCfg.Base)
}

func FinishBranch(branchType, branchName string) {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Could not load .gitkit.yml: %v\n", err)
		os.Exit(1)
	}
	prefixCfg, ok := cfg.Prefixes[branchType]
	if !ok {
		fmt.Printf("❌ Unknown type: %s\n", branchType)
		return
	}
	prefix := prefixCfg.Name
	base := prefixCfg.Base
	var branch string
	if branchName == "" {
		branch = branchName
	} else {
		branch = CurrentBranch()
	}
	branch = RemovePrefix(branch, prefix)
	branch = prefix + branch
	if !BranchExists(branch) {
		fmt.Fprintf(os.Stderr, "❌ branch '%s' does not exist.\n", branch)
		os.Exit(1)
	}
	Checkout(branch)
	Push()
	err = MergeBranchToBase(base, branch)
	if err != nil {
		prompt := promptui.Prompt{
			Label:     "Merge failed. Do you want to create a merge request instead(Y,n)?",
			IsConfirm: true,
		}

		result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		if result != "y" && result != "Y" && result != "" {
			os.Exit(1)
		}
		err = CreatePR(base, branch, "Merge request "+branch, "automatic merge request from "+branch+"to "+base)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not create PR: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("✅ Merge request created. %s --> %s\n", branch, base)
		}
		return
	}

	fmt.Printf("✅ Feature branch '%s' finished and merged into '%s'.\n", branch, base)
}
