package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/samanar/gitkit/config"
)

type GitlabRepo struct {
	Config config.RepoConfig
}

func NewGitlabRepo(rootPath string) GitlabRepo {
	gitlabRepo := GitlabRepo{}
	gitlabRepo.Config = config.NewGitKitRepoConfig(rootPath, config.Gitlab)
	return gitlabRepo
}

func (repo *GitlabRepo) CreateMergeRequest(targetBranch, sourceBranch, title, description string) error {
	apiURL, err := url.JoinPath(repo.Config.Url, "api/v4", "projects",
		url.PathEscape(fmt.Sprintf("%s/%s", repo.Config.Username, repo.Config.RepositoryName)), "merge_requests")
	fmt.Println(apiURL)
	if err != nil {
		return fmt.Errorf("failed to build API URL: %w", err)
	}

	payload := map[string]string{
		"source_branch":        sourceBranch,
		"target_branch":        targetBranch,
		"title":                title,
		"description":          description,
		"remove_source_branch": "false",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Private-Token", repo.Config.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("network or SSL error: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusCreated {
		var result map[string]any
		if err := json.Unmarshal(respBody, &result); err == nil {
			fmt.Printf("✅ Merge Request created: %s\n", result["web_url"])
			return nil
		}
		fmt.Println("✅ Merge Request created successfully.")
		return nil
	}

	switch resp.StatusCode {
	case 400:
		return fmt.Errorf("bad request: %s", string(respBody))
	case 401:
		return fmt.Errorf("unauthorized: invalid or missing token")
	case 403:
		return fmt.Errorf("forbidden: token lacks permissions to create MR")
	case 404:
		return fmt.Errorf("repository not found or you don’t have access: %s", string(respBody))
	case 409:
		return fmt.Errorf("merge request conflict or duplicate: %s", string(respBody))
	default:
		return fmt.Errorf("failed to create MR: [%d] %s", resp.StatusCode, string(respBody))
	}
}
