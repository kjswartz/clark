package cmd

import (
	"bufio"
	"fmt"

	"os"

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

	fmt.Print("What is the name of the owner? ")
	owner, _ := reader.ReadString('\n')
	owner = owner[:len(owner)-1] // Remove the newline character

	fmt.Print("What is the name of the repository? ")
	repo, _ := reader.ReadString('\n')
	repo = repo[:len(repo)-1] // Remove the newline character

	fmt.Printf("Owner: %s, Repository: %s\n", owner, repo)

	err = gh.SubmitPullRequest(owner, repo)
	if err != nil {
		fmt.Println("Error submitting pull request:", err)
		os.Exit(1)
	}
	fmt.Println("Pull request submitted successfully!")
}
