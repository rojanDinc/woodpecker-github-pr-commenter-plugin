package plugin

import "errors"

type Settings struct {
	GithubToken       string
	Repository        string
	Comment           string
	PullRequestNumber int64
	Owner             string
}

func (s *Settings) Validate() error {
	if s.GithubToken == "" {
		return errors.New("GitHub token is required")
	}

	if s.Repository == "" {
		return errors.New("GitHub repository is required")
	}

	if s.PullRequestNumber == 0 {
		return errors.New("pull request number is required")
	}

	if s.Comment == "" {
		return errors.New("comment is required")
	}

	if s.Owner == "" {
		return errors.New("owner is required")
	}

	return nil
}
