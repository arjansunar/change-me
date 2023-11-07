package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var numCommits int
	var allCommits bool

	rootCmd := &cobra.Command{Use: "git-author-rewrite", Short: "Change author details of Git commits while preserving the commit date"}
	rootCmd.PersistentFlags().IntVarP(&numCommits, "num-commits", "n", 1, "Number of commits to change author information for")
	rootCmd.PersistentFlags().BoolVarP(&allCommits, "all", "a", false, "Apply the same author information to all commits")

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Usage: git-author-rewrite [OPTIONS] <new-author> <email>")
			os.Exit(1)
		}

		newAuthor := args[0]
		newEmail := args[1]

		gitLogCommand := "git log --reverse --format='%h'"

		if !allCommits {
			gitLogCommand += " --max-count="
			gitLogCommand += fmt.Sprintf("%d", numCommits)
		}

		output, err := exec.Command("sh", "-c", gitLogCommand).Output()
		if err != nil {
			log.Fatal("Error executing 'git log' command")
		}

		commitHashes := strings.Fields(string(output))

		for _, commitHash := range commitHashes {
			// Retrieve the original commit date
			dateCommand := fmt.Sprintf("git show -s --format='%%at' %s", commitHash)
			dateOutput, err := exec.Command("sh", "-c", dateCommand).Output()
			if err != nil {
				log.Fatal("Error getting commit date")
			}
			originalDate := strings.TrimSpace(string(dateOutput))

			// Create a new environment with the updated author details and original commit date
			env := fmt.Sprintf("GIT_AUTHOR_NAME='%s' GIT_AUTHOR_EMAIL='%s' GIT_AUTHOR_DATE='%s'", newAuthor, newEmail, originalDate)

			// Amend the commit with the new author details and original commit date
			amendCommand := fmt.Sprintf("%s git commit --amend --reset-author --no-edit", env)
			if err := exec.Command("sh", "-c", amendCommand).Run(); err != nil {
				log.Fatalf("Error amending commit %s\n", commitHash)
			}
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
