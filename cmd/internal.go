package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/v53/github"
	"strings"
	"time"
)

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type GitHubRepoPredicate func(client *github.Client, ctx context.Context, repo *github.Repository) (bool, error)
type GitHubRepoAction func(client *github.Client, ctx context.Context, repo *github.Repository) error

func ForEachRepositories(client *github.Client, ctx context.Context, organization string, predicate GitHubRepoPredicate, action GitHubRepoAction) error {
	opt := &github.RepositoryListByOrgOptions{
		Type: "all",
		ListOptions: github.ListOptions{
			PerPage: 500,
		},
	}

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, organization, opt)
		if err != nil {
			return fmt.Errorf("error fetching repositories: %v", err)
		}

		for _, repo := range repos {
			var ok bool
			if ok, err = predicate(client, ctx, repo); err != nil {
				return fmt.Errorf("error evaluating predicate for repo %s: %v", repo.GetName(), err)
			} else if ok {
				if err = action(client, ctx, repo); err != nil {
					return fmt.Errorf("error performing action for repo %s: %v", repo.GetName(), err)
				}
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return nil
}

func isRepositoryEmpty(client *github.Client, ctx context.Context, repo *github.Repository) (bool, error) {
	repoName := repo.GetName()
	defaultBranch := repo.GetDefaultBranch()

	// if repo name ends in .github.io, then skip it, since it's the documentation of the Organization
	if strings.HasSuffix(repoName, ".github.io") {
		return false, nil
	}

	// if repo is .github, skip it, since it's the organization's profile
	if repoName == ".github" {
		return false, nil
	}

	// Get the contents of the default branch
	ref, _, err := client.Git.GetRef(ctx, organization, repoName, "refs/heads/"+defaultBranch)
	if err != nil {
		return false, fmt.Errorf("error getting default branch for repo %s: %v", repoName, err)
	}

	tree, _, err := client.Git.GetTree(ctx, organization, repoName, ref.GetObject().GetSHA(), true)
	if err != nil {
		return false, fmt.Errorf("error getting tree for repo %s: %v", repoName, err)
	}

	fileMap := make(map[string]*github.TreeEntry)
	for _, entry := range tree.Entries {
		if entry.GetType() == "blob" {
			fileMap[entry.GetPath()] = entry
		}
	}
	lastMonth := time.Now().AddDate(0, -1, 0)

	if len(fileMap) > 0 {
		for k := range fileMap {
			l := strings.ToLower(k)
			if !contains([]string{"readme.md", "license", ".gitignore"}, l) {
				return false, nil
			}
		}

		// let's get the time of the last commit, and if that's more than a month ago, we'll consider it a candidate
		commits, _, err := client.Repositories.ListCommits(ctx, organization, repoName, &github.CommitsListOptions{
			SHA: ref.GetObject().GetSHA(),
		})
		if err != nil {
			return false, fmt.Errorf("error getting commits for repo %s: %v", repoName, err)
		}

		if len(commits) > 0 {
			lastCommit := commits[0]
			if lastCommit.GetCommit().GetCommitter().GetDate().After(lastMonth) {
				// last commit is less than a month ago
				return false, nil
			}
		}

		return true, nil
	}
	return true, nil
}

func printRepositoryName(client *github.Client, ctx context.Context, repo *github.Repository) error {
	fmt.Println(repo.GetName())
	return nil
}

func archiveRepository(client *github.Client, ctx context.Context, repo *github.Repository) error {
	// Check if the repository is already archived
	if repo.GetArchived() {
		fmt.Printf("Repository %s is already archived.\n", repo.GetName())
		return nil
	}

	// Prepare the repository object with the Archived field set to true
	repoUpdate := &github.Repository{
		Archived: github.Bool(true),
	}

	// Update the repository to set it as archived
	_, _, err := client.Repositories.Edit(ctx, repo.GetOwner().GetLogin(), repo.GetName(), repoUpdate)
	if err != nil {
		return fmt.Errorf("error archiving repository %s: %v", repo.GetName(), err)
	}

	fmt.Printf("Repository %s has been archived.\n", repo.GetName())

	return nil
}

func confirmAction(message string) bool {
	var response string
	fmt.Printf("%s [y/N]: ", message)
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
