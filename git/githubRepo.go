package git

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

func NewGithubRepo() *GithubConfigStruct {
	cfg := NewGitKitRepoConfig()
	cfg.GetGitHub()
	return &cfg.Github
}

func (cfg *GithubConfigStruct) getClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.AccessToken})
	tc := oauth2.NewClient(ctx, ts)

	if cfg.Url != "" {
		client, err := github.NewEnterpriseClient(cfg.Url, cfg.Url, tc)
		if err != nil {
			log.Fatalf("Failed to create GitHub Enterprise client: %v", err)
		}
		return client
	}

	return github.NewClient(tc)
}

func (cfg *GithubConfigStruct) CreatePR(targetBranch, branch, title, body string) error {
	client := cfg.getClient()
	ctx := context.Background()
	fmt.Printf("title: %s head: %s base: %s body %s", targetBranch, branch, title, body)

	pr := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(branch),
		Base:  github.String(targetBranch),
		Body:  github.String(body),
	}
	fmt.Printf("title: %s head: %s base: %s body %s", pr.Title, pr.Head, pr.Base, pr.Body)

	log.Printf("🔄 Creating PR: %s → %s/%s:%s\n", branch, cfg.Username, cfg.RepositoryName, targetBranch)

	newPR, resp, err := client.PullRequests.Create(ctx, cfg.Username, cfg.RepositoryName, pr)
	if err != nil {
		if resp != nil && resp.StatusCode == 422 {
			return fmt.Errorf("PR creation failed (possibly duplicate PR or invalid branch): %w", err)
		}
		if strings.Contains(err.Error(), "401") {
			return fmt.Errorf("authentication failed: check your GitHub token in config.json")
		}
		if strings.Contains(err.Error(), "404") {
			return fmt.Errorf("repository not found or no access to %s/%s", cfg.Username, cfg.RepositoryName)
		}
		return fmt.Errorf("failed to create PR: %w", err)
	}

	fmt.Printf("✅ Pull Request created: %s\n", newPR.GetHTMLURL())
	return nil
}
