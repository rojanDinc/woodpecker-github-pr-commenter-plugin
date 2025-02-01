package plugin

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestPlugin(t *testing.T) {
	cases := []struct {
		name          string
		settings      *Settings
		resp          *http.Response
		expectedError error
	}{
		{
			name: "happy path",
			settings: &Settings{
				GithubToken:       "token",
				Repository:        "repo",
				Comment:           "comment",
				PullRequestNumber: 1,
				Owner:             "owner",
			},
			resp: &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewBufferString("")),
				Header:     make(http.Header),
			},
			expectedError: nil,
		},
		{
			name: "failed to create comment",
			settings: &Settings{
				GithubToken:       "token",
				Repository:        "repo",
				Comment:           "comment",
				PullRequestNumber: 1,
				Owner:             "owner",
			},
			resp: &http.Response{
				StatusCode: 403,
				Body:       io.NopCloser(bytes.NewBufferString(``)),
				Header:     make(http.Header),
			},
			expectedError: errors.New("failed to create comment got unexpected status code"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewTestClient(func(req *http.Request) *http.Response {
				return tc.resp
			})
			p := NewPlugin("", client, tc.settings)

			err := p.Execute(context.Background())

			assert.Equal(t, err, tc.expectedError)
		})
	}
}
