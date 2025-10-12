package git

import "github.com/samanar/gitkit/config"

type GitCmd struct {
	Config config.GitKitConfig
}

func NewGitCmdWithoutConfig() GitCmd {
	gitCmd := GitCmd{}
	gitCmd.Config = config.GitKitConfig{}
	return gitCmd
}

func NewGitCommandWithConfig(overwriteConfig bool) GitCmd {
	gitCmd := GitCmd{}
	rootPath, err := gitCmd.RootDir()
	if err != nil {
		panic(err)
	}
	gitCmd.Config = config.NewGitConfig(rootPath, overwriteConfig)
	return gitCmd
}
