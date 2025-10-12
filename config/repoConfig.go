package config

import (
	"fmt"
	"gitkit/gitignore"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type RepoConfig struct {
	Username       string `yaml:"username"`
	AccessToken    string `yaml:"accessToken"`
	Url            string `yaml:"url"`
	RepositoryName string `yaml:"repositoryName"`
}

type GitRepoType int

const (
	Github GitRepoType = iota
	Gitlab
)

const GITKIT_REPO_CONFIG_FILE = ".gitkit_private.yml"

func NewGitKitRepoConfig(rootPath string, repoType GitRepoType) RepoConfig {
	config := RepoConfig{}
	err := config.Load(rootPath, repoType)
	if err != nil {
		os.Exit(1)
	}
	return config
}

func (cfg *RepoConfig) Load(rootPath string, repoType GitRepoType) error {
	path := filepath.Join(rootPath, GITKIT_REPO_CONFIG_FILE)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return err
	}
	return nil
}

func (cfg *RepoConfig) Create(rootPath string, repoType GitRepoType) error {
	defaultUrl := "https://api.github.com/"
	repoName := "Github"
	if repoType == Gitlab {
		repoName = "Gitlab"
		defaultUrl = "https://gitlab.com/"
	}
	fmt.Println("Let's set up your " + repoName + " configuration:")
	cfg.Username = Ask(repoName+" Username", "")
	cfg.AccessToken = Ask("Access Token", "")
	cfg.RepositoryName = Ask("Repository name", "")
	cfg.Url = Ask(repoName+" URL", defaultUrl)
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	path := filepath.Join(rootPath, GITKIT_REPO_CONFIG_FILE)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	fmt.Printf("✅ %s created successfully!\n", GITKIT_REPO_CONFIG_FILE)
	gitignore.AddToGitignore(rootPath, GITKIT_REPO_CONFIG_FILE)
	return nil
}

// func CreatePR(targetBranch, branch, title, body string) error {
// 	cfg, err := LoadConfig()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "❌ Could not load .gitkit.yml: %v\n", err)
// 		os.Exit(1)
// 	}

// 	switch cfg.Repo {
// 	case "Github":
// 		github := NewGithubRepo()
// 		err = github.CreatePR(targetBranch, branch, title, body)
// 		return err
// 		// if err != nil {
// 		// 	fmt.Fprintf(os.Stderr, "❌ Could not create PR: %v\n", err)
// 		// 	os.Exit(1)
// 		// }
// 	case "Gitlab":
// 		github := NewGitlabRepo()
// 		err = github.CreateMergeRequest(targetBranch, branch, title, body)
// 		return err
// 	default:
// 		return fmt.Errorf("unsupported repo type: %s", cfg.Repo)
// 	}
// }
