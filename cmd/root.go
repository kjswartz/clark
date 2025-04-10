package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kjswartz/clark/pkg/githubclient"
)

func RootCmd() {
	// Initialize the GitHub client
	gh, err := githubclient.NewGitHubClient()
	if err != nil {
		fmt.Println("Error creating GitHub client:", err)
		os.Exit(1)
	}

	// Collect data from the user
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("What is the name for the owner? ")
	owner, _ := reader.ReadString('\n')
	owner = owner[:len(owner)-1] // Remove the newline character

	fmt.Print("What are the names of the repositories? ")
	repoInput, _ := reader.ReadString('\n')
	repoInput = repoInput[:len(repoInput)-1] // Remove the newline character

	// Split the input by "," and trim whitespace
	repos := []string{}
	for _, repo := range strings.Split(repoInput, ",") {
		repos = append(repos, strings.TrimSpace(repo))
	}

	fmt.Printf("Owner: %s, Repositories: %v\n", owner, repos)

	for _, repo := range repos {
		fmt.Printf("Submitting pull request for %s/%s...\n", owner, repo)
		// Submit the pull request
		err = gh.SubmitPullRequest(owner, repo)
		if err != nil {
			fmt.Printf("Error submitting pull request for %s/%s: %v", owner, repo, err)
			os.Exit(1)
		}

		fmt.Printf("Pull request submitted successfully for %s/%s...\n", owner, repo)
	}
}
