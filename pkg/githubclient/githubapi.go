package githubclient

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/go-github/v71/github"
	"golang.org/x/oauth2"
)

// GitHubAPI is a wrapper around the GitHub API.
type GitHubAPI struct {
	client *github.Client
}

// newGitHubAPI creates a new client for interacting with GitHub API.
func newGitHubAPI() (*GitHubAPI, error) {
	token := os.Getenv("TOKEN")
	if token == "" {
		return nil, errors.New("TOKEN environment variable not set")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return &GitHubAPI{client: client}, nil
}

func (gh GitHubAPI) GetBranch(owner, repo, name string) (*github.Branch, error) {
	branch, _, err := gh.client.Repositories.GetBranch(context.Background(), owner, repo, name, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}

	return branch, nil
}

// SubmitIssue submits an issue to GitHub.
func (gh GitHubAPI) SubmitPullRequest(owner, repo, head string) error {
	title := "Automated Pull Request"
	body := "This is an automated pull request to add files to your repo from clark."
	base := "main"

	// Create the pull request
	_, _, err := gh.client.PullRequests.Create(context.Background(), owner, repo, &github.NewPullRequest{
		Title: &title,
		Body:  &body,
		Base:  &base,
		Head:  &head,
	})

	return err
}
