package git

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type GithubConfigStruct struct {
	Username       string `yaml:"username"`
	AccessToken    string `yaml:"accessToken"`
	Url            string `yaml:"url"`
	RepositoryName string `yaml:"repositoryName"`
}

type GitlabConfigStruct struct {
	Username    string `yaml:"username"`
	AccessToken string `yaml:"accessToken"`
	Url         string `yaml:"url"`
}

type GitKitRepoConfig struct {
	Github GithubConfigStruct `yaml: "github"`
	Gitlab GitlabConfigStruct `yaml: "gitlab"`
}

type GitConfigType int

var PrivateConfigFile string = ".gitkit_private.yml"

const (
	Github GitConfigType = iota
	Gitlab
)

func NewGitKitRepoConfig() *GitKitRepoConfig {
	return &GitKitRepoConfig{
		Github: GithubConfigStruct{},
		Gitlab: GitlabConfigStruct{},
	}
}

func (cfg *GitKitRepoConfig) Load() (*GitKitRepoConfig, error) {
	root, err := RootDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(root, PrivateConfigFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *GitKitRepoConfig) CreateGithub() (GithubConfigStruct, error) {
	fmt.Println("Let's set up your github configuration:")
	cfg.Github.Username = Ask("Github Username", "")
	cfg.Github.AccessToken = Ask("Github Access Token", "")
	cfg.Github.RepositoryName = Ask("Repository name", "")
	cfg.Github.Url = Ask("Github URL", "https://api.github.com/")
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return GithubConfigStruct{}, err
	}
	err = cfg.CreateConfigFile(data)
	if err != nil {
		return GithubConfigStruct{}, err
	}
	return cfg.Github, nil
}

func (cfg *GitKitRepoConfig) CreateConfigFile(data []byte) error {
	root, err := RootDir()
	if err != nil {
		return err
	}
	path := filepath.Join(root, PrivateConfigFile)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}

func (cfg *GitKitRepoConfig) GetGitHub() (GithubConfigStruct, error) {
	config, err := cfg.Load()
	if err != nil {
		return cfg.CreateGithub()
	}
	gitHubConfig := config.Github
	if gitHubConfig.Username == "" || gitHubConfig.AccessToken == "" || gitHubConfig.Url == "" {
		return cfg.CreateGithub()
	}
	return gitHubConfig, nil
}

func (cfg *GitKitRepoConfig) GetGitlab() (*GitlabConfigStruct, error) {
	return &GitlabConfigStruct{}, fmt.Errorf("method not implemented yet")
}

func (cfg *GitKitRepoConfig) Get(configType GitConfigType) (interface{}, error) {
	switch configType {
	case Github:
		gitHubConfig, err := cfg.GetGitHub()
		if err != nil {
			return nil, err
		}
		return gitHubConfig, nil
	case Gitlab:
		gitLabConfig, err := cfg.GetGitlab()
		if err != nil {
			return nil, err
		}
		return gitLabConfig, nil
	default:
		return nil, fmt.Errorf("unsopported configType %v", configType)
	}
}
