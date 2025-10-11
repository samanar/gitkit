package git

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v3"
)

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

	// AutoSync bool   `yaml:"autoSync"`
}

func CreateConfig() error {
	repoRoot, err := RootDir()
	if err != nil {
		return fmt.Errorf("not a git repository: %v", err)
	}

	cfgPath := filepath.Join(repoRoot, ".gitkit.yml")

	if _, err := os.Stat(cfgPath); err == nil {
		fmt.Printf("⚠️  Config file already exists at %s\n", cfgPath)
		fmt.Print("Do you want to overwrite it? (y/N): ")
		confirm := strings.ToLower(strings.TrimSpace(ReadLine()))
		if confirm != "y" && confirm != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	cfg := GitKitConfig{}

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
	cfg.Repo = result
	cfg.Branches.Main = Ask("Main branch name", "main")
	cfg.Branches.Develop = Ask("Develop branch name", "develop")
	cfg.Prefixes = make(map[string]struct {
		Name string `yaml:"name"`
		Base string `yaml:"base"`
	})
	cfg.Prefixes["feature"] = struct {
		Name string `yaml:"name"`
		Base string `yaml:"base"`
	}{
		Name: Ask("Feature prefix", "feature/"),
		Base: cfg.Branches.Develop,
	}
	cfg.Prefixes["bugFix"] = struct {
		Name string `yaml:"name"`
		Base string `yaml:"base"`
	}{
		Name: Ask("Bugfix prefix", "bugfix/"),
		Base: cfg.Branches.Develop,
	}
	cfg.Prefixes["hotFix"] = struct {
		Name string `yaml:"name"`
		Base string `yaml:"base"`
	}{
		Name: Ask("Hotfix prefix", "hotfix/"),
		Base: cfg.Branches.Main,
	}
	cfg.Prefixes["release"] = struct {
		Name string `yaml:"name"`
		Base string `yaml:"base"`
	}{
		Name: Ask("Release prefix", "release/"),
		Base: cfg.Branches.Develop,
	}
	cfg.Remote = Ask("Remote name", "origin")

	// Save YAML file
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile(cfgPath, data, 0644); err != nil {
		return err
	}

	fmt.Printf("✅ Configuration saved to %s\n", cfgPath)
	return nil
}

func LoadConfig() (*GitKitConfig, error) {
	root, err := RootDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(root, ".gitkit.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		err = CreateConfig()
		if err != nil {
			return nil, err
		}
		return LoadConfig()
	}

	var cfg GitKitConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func ReadLine() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func Ask(question, def string) string {
	fmt.Printf("%s [%s]: ", question, def)
	answer := ReadLine()
	if answer == "" {
		return def
	}
	return answer
}
