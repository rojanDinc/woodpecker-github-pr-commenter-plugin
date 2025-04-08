package plugin

type Settings struct {
	GithubToken       string
	Repository        string
	Comment           string
	PullRequestNumber int64
	Owner             string
	LogLevel          string
}
