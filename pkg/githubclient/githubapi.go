package githubclient

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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

// SubmitIssue submits an issue to GitHub.
func (gh GitHubAPI) SubmitPullRequest(owner, repo string) error {
	fmt.Println("Submitting pull request to GitHub...")

	title := "Automated Pull Request"
	body := "This is an automated pull request to add files to your repo from clark."
	base := "main"
	head := "clark-install-branch"

	// Clone the cold-harbor repository
	tempDir, err := os.MkdirTemp("", "cold-harbor-clone")
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

	// Check if the head branch exists
	_, _, err = gh.client.Repositories.GetBranch(context.Background(), owner, repo, head, 0)
	if err != nil {
		fmt.Printf("Branch '%s' does not exist. Creating it...\n", head)

		// Create the branch
		baseBranch, _, err := gh.client.Repositories.GetBranch(context.Background(), owner, repo, base, 0)
		if err != nil {
			return fmt.Errorf("failed to get base branch '%s': %w", base, err)
		}

		newBranchRef := fmt.Sprintf("refs/heads/%s", head)
		ref := &github.Reference{
			Ref: &newBranchRef,
			Object: &github.GitObject{
				SHA: baseBranch.Commit.SHA,
			},
		}

		_, _, err = gh.client.Git.CreateRef(context.Background(), owner, repo, ref)
		if err != nil {
			return fmt.Errorf("failed to create branch '%s': %w", head, err)
		}
	}

	// Create the pull request
	_, _, err = gh.client.PullRequests.Create(context.Background(), owner, repo, &github.NewPullRequest{
		Title: &title,
		Body:  &body,
		Base:  &base,
		Head:  &head,
	})

	return err
}
