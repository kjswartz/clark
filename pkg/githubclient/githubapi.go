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
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, errors.New("GITHUB_TOKEN environment variable not set")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return &GitHubAPI{client: client}, nil
}

// SubmitIssue submits an issue to GitHub.
func (gh GitHubAPI) SubmitPullRequest(owner, repo string) error {
	fmt.Println("Submitting pull request to GitHub...")
	return nil
	// _, _, err := gh.client.Issues.Create(context.Background(), owner, repo, &github.IssueRequest{
	// 	Title: "&title",
	// 	Body:  "&body",
	// })

	// return err
}
