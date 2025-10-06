package config

import (
	"bufio"
	"fmt"
	"gitkit/git"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type GitKitConfig struct {
	Branches struct {
		Main    string `yaml:"main"`
		Develop string `yaml:"develop"`
	} `yaml:"branches"`
	Prefixes struct {
		Feature string `yaml:"feature"`
		Hotfix  string `yaml:"hotfix"`
		Release string `yaml:"release"`
	} `yaml:"prefixes"`
	Remote string `yaml:"remote"`
	// AutoSync bool   `yaml:"autoSync"`
}

func LoadConfig() (*GitKitConfig, error) {
	root, err := git.RootDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(root, ".gitflow.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
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
