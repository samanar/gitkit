package git

import (
	"github.com/samanar/gitkit/config"
)

type GitCmd struct {
	Config config.GitKitConfig
}

func NewGitCmdWithoutConfig() GitCmd {
	gitCmd := GitCmd{}
	rootPath, err := gitCmd.RootDir()
	if err != nil {
		gitCmd.Config = config.GitKitConfig{}
	} else {
		gitCmd.Config = config.NewGitConfig(rootPath, false)
	}

	return gitCmd
}

func NewGitCommandWithConfig() GitCmd {
	gitCmd := GitCmd{}
	rootPath, err := gitCmd.RootDir()
	if err != nil {
		panic(err)
	}
	gitCmd.Config = config.NewGitConfig(rootPath, true)
	return gitCmd
}

func NewGitCommandWithForceRewriteConfig() GitCmd {
	gitCmd := GitCmd{}
	rootPath, err := gitCmd.RootDir()
	if err != nil {
		panic(err)
	}
	gitCmd.Config = config.NewGitConfigForced(rootPath)
	return gitCmd
}
