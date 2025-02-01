package main

import (
	"context"
	"net/http"
	"os"

	"github.com/rojanDinc/woodpecker-github-pr-commenter-plugin/internal/command"
	"github.com/urfave/cli/v3"
)

func main() {
	createCmd := command.NewCreate(http.DefaultClient)
	cmd := &cli.Command{
		Commands: []*cli.Command{
			createCmd.Command(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
