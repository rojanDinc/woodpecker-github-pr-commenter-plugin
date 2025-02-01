package command

import (
	"context"
	"net/http"

	"github.com/rojanDinc/woodpecker-github-pr-commenter-plugin/internal/plugin"
	"github.com/urfave/cli/v3"
)

type Create struct {
	httpClient *http.Client
	settings   *plugin.Settings
}

func NewCreate(httpClient *http.Client) *Create {
	return &Create{
		httpClient: httpClient,
		settings:   &plugin.Settings{},
	}
}

func (c *Create) Command() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "create a new PR comment",
		Flags: c.flags(),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			plugin := plugin.NewPlugin("https://api.github.com", c.httpClient, c.settings)

			return plugin.Execute(ctx)
		},
	}

}

func (c *Create) flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "token",
			Usage:       "GitHub bearer token to be used for authentication",
			Sources:     cli.EnvVars("PLUGIN_GITHUB_TOKEN"),
			Destination: &c.settings.GithubToken,
		},
		&cli.StringFlag{
			Name:        "repo",
			Usage:       "GitHub repository",
			Sources:     cli.EnvVars("CI_REPO"),
			Destination: &c.settings.Repository,
		},
		&cli.IntFlag{
			Name:        "pr-number",
			Usage:       "Pull request number",
			Sources:     cli.EnvVars("CI_COMMIT_PULL_REQUEST"),
			Destination: &c.settings.PullRequestNumber,
		},
		&cli.StringFlag{
			Name:        "comment",
			Usage:       "Comment to be added to pull request",
			Sources:     cli.EnvVars("PLUGIN_COMMENT"),
			Destination: &c.settings.Comment,
		},
	}
}
