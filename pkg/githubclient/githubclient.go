// Package githubclient provides all interactions with the GitHub API.
package githubclient

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/go-github/v71/github"
)

// Factory defines interface to use the GitHub API.
type Factory interface {
	SubmitPullRequest(owner, repo, head string) error
	GetBranch(owner, repo, name string) (*github.Branch, error)
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

// SubmitPullRequest submits a pull-request to GitHub.
func (gh GitHubClient) SubmitPullRequest(owner, repo string) error {
	head := "clark-install-branch"

	// Check if the head branch exists
	branch, _ := gh.client.GetBranch(owner, repo, head)
	if branch != nil {
		return fmt.Errorf("Branch '%s' already exists. Skiping...\n", head)
	}

	// Clone the repo
	fmt.Println("Cloning repository and adding files...")
	err := cloneRepo(owner, repo, head)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return gh.client.SubmitPullRequest(owner, repo, head)
}

func cloneRepo(owner, repo, head string) error {
	// Clone the repo
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("%s-clone", repo))
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	repoURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
	cmd := exec.Command("git", "clone", repoURL, tempDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Copy resources to the .github directory in the cloned repository
	sourcePath := "resources"
	destinationPath := filepath.Join(tempDir, ".github")
	err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		destination := filepath.Join(destinationPath, relPath)
		destinationDir := filepath.Dir(destination)

		err = os.MkdirAll(destinationDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destinationDir, err)
		}

		input, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		err = os.WriteFile(destination, input, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", destination, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to copy resources: %w", err)
	}

	// Check out or create the clark-install-branch
	cmd = exec.Command("git", "-C", tempDir, "checkout", "-B", head)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create or switch to branch '%s': %w", head, err)
	}

	// Stage, commit, and push changes
	cmd = exec.Command("git", "-C", tempDir, "add", ".")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	cmd = exec.Command("git", "-C", tempDir, "commit", "-m", "Add resources to .github directory")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	cmd = exec.Command("git", "-C", tempDir, "push", "-u", "origin", head)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push changes to branch '%s': %w", head, err)
	}

	return nil
}
