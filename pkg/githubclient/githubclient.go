// Package githubclient provides all interactions with the GitHub API.
package githubclient

// Factory defines interface to use the GitHub API.
type Factory interface {
	SubmitPullRequest(owner, repo string) error
}

// GitHubClient - The GitHubClient struct has a client that uses the GitHub REST API.
type GitHubClient struct {
	client Factory
}

// NewGitHubClient creates a new Github client. This allows for mocking the GitHub API in tests.
func NewGitHubClient() (gitHubClient *GitHubClient, err error) {
	gh, err := newGitHubAPI()
	if err != nil {
		return nil, err
	}
	return &GitHubClient{client: gh}, nil
}

// SubmitIssue submits an issue to GitHub.
func (gh GitHubClient) SubmitPullRequest(owner, repo string) error {
	return gh.client.SubmitPullRequest(owner, repo)
}
