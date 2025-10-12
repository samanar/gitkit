package repo

import (
	"fmt"

	"github.com/samanar/gitkit/config"
)

func MergeRequest(rootPath string, repoType config.GitRepoType, title, description, sourceBranch, targetBranch string) error {
	switch repoType {
	case config.Gitlab:
		gitlabRepo := NewGitlabRepo(rootPath)
		return gitlabRepo.CreateMergeRequest(targetBranch, sourceBranch, title, description)
	case config.Github:
		githubRepo := NewGithubRepo(rootPath)
		return githubRepo.CreatePR(targetBranch, sourceBranch, title, description)
	default:
		return fmt.Errorf("unsupported repo type: %v", repoType)
	}
}
