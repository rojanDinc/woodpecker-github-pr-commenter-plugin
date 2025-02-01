package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type CreateCommentRequest struct {
	Body string `json:"body"`
}

type Plugin struct {
	baseURL    string
	httpClient *http.Client
	Settings   *Settings
}

func NewPlugin(ghBaseURL string, httpClient *http.Client, settings *Settings) *Plugin {
	return &Plugin{
		baseURL:    ghBaseURL,
		httpClient: httpClient,
		Settings:   settings,
	}

}

func (p *Plugin) Execute(ctx context.Context) error {
	if err := p.Settings.Validate(); err != nil {
		return err
	}

	data := CreateCommentRequest{
		Body: p.Settings.Comment,
	}
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments", p.baseURL, p.Settings.Owner, p.Settings.Repository, p.Settings.PullRequestNumber), buf)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.Settings.GithubToken))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return errors.New("failed to create comment got unexpected status code")
	}

	return nil
}
