package actions

import (
	"log"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

// HelmGenerateE2E is the cli handler for generating a helm parameters file for deis-e2e
func HelmGenerateE2E(ghClient *github.Client) func(*cli.Context) {
	return func(c *cli.Context) {
		const (
			repoName = "workflow-e2e"
			chartDir = "workflow-dev-e2e"
		)
		params := genParamsComponentAttrs{
			Org:        c.GlobalString(OrgFlag),
			PullPolicy: c.GlobalString(PullPolicyFlag),
			Tag:        c.GlobalString(TagFlag),
		}
		paramsComponentMap := createParamsComponentMap()
		paramsComponentMap[repoName] = params

		if params.Tag == "" {
			reposAndShas, err := getShas(ghClient, []string{repoName}, shortShaTransform)
			if err != nil {
				log.Fatalf("No tag given and couldn't fetch sha from GitHub (%s)", err)
			} else if len(reposAndShas) < 1 {
				log.Fatalf("No tag given and no sha returned from GitHub for deis/%s", repoName)
			}
			params.Tag = "git-" + reposAndShas[0].sha
			paramsComponentMap[repoName] = params
		}
		shouldStage := c.GlobalBool(StageFlag)
		stagingDir := filepath.Join(stagingPath, chartDir)
		if err := generateParams(shouldStage, ourFS, stagingDir, paramsComponentMap); err != nil {
			log.Fatalf("Error outputting the workflow values file (%s)", err)
		}
	}
}
