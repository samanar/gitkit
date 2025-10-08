package cmd

import (
	"fmt"
	"gitkit/git"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoRoot, err := git.RootDir()
		if err != nil {
			return fmt.Errorf("not a git repository: %v", err)
		}

		cfgPath := filepath.Join(repoRoot, ".gitkit.yml")

		if _, err := os.Stat(cfgPath); err == nil {
			fmt.Printf("⚠️  Config file already exists at %s\n", cfgPath)
			fmt.Print("Do you want to overwrite it? (y/N): ")
			confirm := strings.ToLower(strings.TrimSpace(git.ReadLine()))
			if confirm != "y" && confirm != "yes" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		cfg := git.GitKitConfig{}

		// Ask user inputs with defaults
		fmt.Println("Let's set up your gitkit configuration:")
		cfg.Branches.Main = git.Ask("Main branch name", "main")
		cfg.Branches.Develop = git.Ask("Develop branch name", "develop")
		cfg.Prefixes = make(map[string]struct {
			Name string `yaml:"name"`
			Base string `yaml:"base"`
		})
		cfg.Prefixes["feature"] = struct {
			Name string `yaml:"name"`
			Base string `yaml:"base"`
		}{
			Name: git.Ask("Feature prefix", "feature/"),
			Base: cfg.Branches.Develop,
		}
		cfg.Prefixes["bugFix"] = struct {
			Name string `yaml:"name"`
			Base string `yaml:"base"`
		}{
			Name: git.Ask("Bugfix prefix", "bugfix/"),
			Base: cfg.Branches.Develop,
		}
		cfg.Prefixes["hotFix"] = struct {
			Name string `yaml:"name"`
			Base string `yaml:"base"`
		}{
			Name: git.Ask("Hotfix prefix", "hotfix/"),
			Base: cfg.Branches.Main,
		}
		cfg.Prefixes["release"] = struct {
			Name string `yaml:"name"`
			Base string `yaml:"base"`
		}{
			Name: git.Ask("Release prefix", "release/"),
			Base: cfg.Branches.Develop,
		}
		cfg.Remote = git.Ask("Remote name", "origin")

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
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
