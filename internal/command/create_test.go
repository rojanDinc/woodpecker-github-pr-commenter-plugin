package command

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestCreateCommandSuccess(t *testing.T) {
	vars := map[string]string{
		"PLUGIN_GITHUB_TOKEN":    "token",
		"CI_REPO":                "repo",
		"CI_COMMIT_PULL_REQUEST": "1",
		"PLUGIN_COMMENT":         "comment",
		"CI_REPO_OWNER":          "owner",
		"PLUGIN_LOG_LEVEL":       "debug",
	}
	for k, v := range vars {
		t.Setenv(k, v)
	}
	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(bytes.NewBufferString("")),
			Header:     make(http.Header),
		}
	})

	createCmd := NewCreate(client)

	cmd := cli.Command{
		Commands: []*cli.Command{
			createCmd.Command(),
		},
	}

	err := cmd.Run(context.Background(), []string{"", "create"})

	assert.NoError(t, err)
	assert.True(t, slog.Default().Handler().Enabled(context.Background(), slog.LevelDebug))
}

func TestCreateCommandFail(t *testing.T) {
	cases := []struct {
		name     string
		vars     map[string]string
		expected error
	}{
		{
			name: "missing token",
			vars: map[string]string{
				"CI_REPO":                "repo",
				"CI_COMMIT_PULL_REQUEST": "1",
				"PLUGIN_COMMENT":         "comment",
				"CI_REPO_OWNER":          "owner",
			},
			expected: errors.New("GitHub token is required"),
		},
		{
			name: "missing repo",
			vars: map[string]string{
				"CI_COMMIT_PULL_REQUEST": "1",
				"PLUGIN_COMMENT":         "comment",
				"PLUGIN_GITHUB_TOKEN":    "token",
				"CI_REPO_OWNER":          "owner",
			},
			expected: errors.New("GitHub repository is required"),
		},
		{
			name: "missing pull request number",
			vars: map[string]string{
				"CI_REPO":             "repo",
				"PLUGIN_COMMENT":      "comment",
				"PLUGIN_GITHUB_TOKEN": "token",
				"CI_REPO_OWNER":       "owner",
			},
			expected: errors.New("pull request number is required"),
		},
		{
			name: "missing comment",
			vars: map[string]string{
				"CI_REPO":                "repo",
				"CI_COMMIT_PULL_REQUEST": "1",
				"PLUGIN_GITHUB_TOKEN":    "token",
				"CI_REPO_OWNER":          "owner",
			},
			expected: errors.New("comment is required"),
		},
		{
			name: "missing owner",
			vars: map[string]string{
				"CI_REPO":                "repo",
				"CI_COMMIT_PULL_REQUEST": "1",
				"PLUGIN_COMMENT":         "comment",
				"PLUGIN_GITHUB_TOKEN":    "token",
			},
			expected: errors.New("owner is required"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for k, v := range c.vars {
				t.Setenv(k, v)
			}

			client := NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 201,
					Body:       io.NopCloser(bytes.NewBufferString("")),
					Header:     make(http.Header),
				}
			})

			createCmd := NewCreate(client)

			cmd := cli.Command{
				Commands: []*cli.Command{
					createCmd.Command(),
				},
			}

			err := cmd.Run(context.Background(), []string{"", "create"})

			assert.Equal(t, c.expected, err)
		})
	}
}
