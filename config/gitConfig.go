package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v3"
)

const GITKIT_CONFIG_FILE = ".gitkit.yaml"

type GitKitConfig struct {
	Repo     string `yaml:"repo"`
	Branches struct {
		Main    string `yaml:"main"`
		Develop string `yaml:"develop"`
	} `yaml:"branches"`
	Prefixes map[string]struct {
		Name string `yaml:"name"`
		Base string `yaml:"base"`
	} `yaml:"prefixes"`
	Remote string `yaml:"remote"`
}

func NewGitConfig(rootPath string, replaceConfig bool) GitKitConfig {
	gitConfig := GitKitConfig{}
	if replaceConfig {
		gitConfig.Create(rootPath)
	}
	gitConfig.Load(rootPath, replaceConfig)
	return gitConfig
}

func (cfg *GitKitConfig) Exists(rootPath string) bool {
	configPath := filepath.Join(rootPath, GITKIT_CONFIG_FILE)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (cfg *GitKitConfig) Load(rootPath string, replaceConfig bool) error {
	if !cfg.Exists(rootPath) {
		cfg.Create(rootPath)
	}
	configPath := filepath.Join(rootPath, GITKIT_CONFIG_FILE)
	data, err := os.ReadFile(configPath)
	if err != nil {
		err = cfg.Create(rootPath)
		if err != nil {
			return err
		}
	}

	var gitkitCfg GitKitConfig
	if err := yaml.Unmarshal(data, &gitkitCfg); err != nil {
		return err
	}
	cfg = &gitkitCfg
	return nil
}

func (cfg *GitKitConfig) Create(rootPath string) error {
	cfgPath := filepath.Join(rootPath, GITKIT_CONFIG_FILE)
	if _, err := os.Stat(cfgPath); err == nil {
		fmt.Printf("⚠️  Config file already exists at %s\n", cfgPath)
		fmt.Print("Do you want to overwrite it? (y/N): ")
		confirm := strings.ToLower(strings.TrimSpace(ReadLine()))
		if confirm != "y" && confirm != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	gitkitConfig := GitKitConfig{}

	// Ask user inputs with defaults
	fmt.Println("Let's set up your gitkit configuration:")
	prompt := promptui.Select{
		Label: "What is your repository hosting service?",
		Items: []string{"Github", "Gitlab"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		return fmt.Errorf("error getting your hosting service: %v", err)
	}
	gitkitConfig.Repo = result
	gitkitConfig.Branches.Main = Ask("Main branch name", "main")
	gitkitConfig.Branches.Develop = Ask("Develop branch name", "develop")
	gitkitConfig.Prefixes = make(map[string]struct {
		Name string `yaml:"name"`
		Base string `yaml:"base"`
	})
	gitkitConfig.Prefixes["feature"] = struct {
		Name string `yaml:"name"`
		Base string `yaml:"base"`
	}{
		Name: Ask("Feature prefix", "feature/"),
		Base: gitkitConfig.Branches.Develop,
	}
	gitkitConfig.Prefixes["bugFix"] = struct {
		Name string `yaml:"name"`
		Base string `yaml:"base"`
	}{
		Name: Ask("Bugfix prefix", "bugfix/"),
		Base: gitkitConfig.Branches.Develop,
	}
	gitkitConfig.Prefixes["hotFix"] = struct {
		Name string `yaml:"name"`
		Base string `yaml:"base"`
	}{
		Name: Ask("Hotfix prefix", "hotfix/"),
		Base: gitkitConfig.Branches.Main,
	}
	gitkitConfig.Prefixes["release"] = struct {
		Name string `yaml:"name"`
		Base string `yaml:"base"`
	}{
		Name: Ask("Release prefix", "release/"),
		Base: gitkitConfig.Branches.Develop,
	}
	gitkitConfig.Remote = Ask("Remote name", "origin")

	// Save YAML file
	data, err := yaml.Marshal(&gitkitConfig)
	if err != nil {
		return err
	}

	if err := os.WriteFile(cfgPath, data, 0644); err != nil {
		return err
	}

	fmt.Printf("✅ Configuration saved to %s\n", cfgPath)
	return nil
}
