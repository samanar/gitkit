package git

import (
	"fmt"
	"gitkit/config"
	"gitkit/repo"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func (g *GitCmd) CurrentBranch() string {
	output := g.RunMust("branch", "--show-current")
	return output[:len(output)-1]
}

func (g *GitCmd) Branches() []string {
	output := g.RunMust("branch", "--list")
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

func (g *GitCmd) BranchesWithPrefix(prefix string) []string {
	branches := g.Branches()
	var matches []string

	for _, branch := range branches {
		if strings.HasPrefix(branch, prefix) {
			matches = append(matches, branch)
		}
	}
	return matches
}

func (g *GitCmd) DeleteBranch(branch string, force bool) {
	flag := "-d"
	if force {
		flag = "-D"
	}
	g.RunMust("branch", flag, branch)
}

func (g *GitCmd) DeleteBranchSafe(branch string) {
	g.DeleteBranch(branch, false)
}

func (g *GitCmd) Checkout(branch string) {
	g.RunMust("checkout", branch)
}

func (g *GitCmd) RemovePrefix(branch, prefix string) string {
	if strings.HasPrefix(branch, prefix) {
		return branch[len(prefix):]
	}
	return branch
}

func (g *GitCmd) BranchExists(branch string) bool {
	branches := g.Branches()
	for _, b := range branches {
		if b == branch {
			return true
		}
	}
	return false
}

func (g *GitCmd) StartBranch(branchType, branchName string) {
	prefixCfg, ok := g.Config.Prefixes[branchType]
	if !ok {
		fmt.Printf("❌ Unknown type: %s\n", branchType)
		return
	}
	branchName = g.RemovePrefix(branchName, prefixCfg.Name)

	g.Sync(prefixCfg.Base)

	// // Use reusable method to create the branch
	if err := g.CreatePrefixedBranch(prefixCfg.Base, prefixCfg.Name, branchName); err != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Feature branch '%s%s' started from '%s'.\n", prefixCfg.Name, branchName, prefixCfg.Base)
}

func (g *GitCmd) FinishBranch(branchType, branchName string) {
	prefixCfg, ok := g.Config.Prefixes[branchType]
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
		branch = g.CurrentBranch()
	}
	branch = g.RemovePrefix(branch, prefix)
	branch = prefix + branch
	if !g.BranchExists(branch) {
		fmt.Fprintf(os.Stderr, "❌ branch '%s' does not exist.\n", branch)
		os.Exit(1)
	}

	err := g.MergeBranchToBase(base, branch)
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
		g.Checkout(branch)
		fmt.Println("current branch is", g.CurrentBranch())
		g.Push()
		repoType := config.Github
		if g.Config.Repo == "Gitlab" {
			repoType = config.Gitlab
		}
		root, err := g.RootDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Could not find git root: %v\n", err)
			os.Exit(1)
		}
		err = repo.MergeRequest(root, repoType, "Merge request "+branch, "automatic merge request from "+branch+"to "+base, branch, base)
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
