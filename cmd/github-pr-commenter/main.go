package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/rojanDinc/woodpecker-github-pr-commenter-plugin/internal/command"
	"github.com/urfave/cli/v3"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	createCmd := command.NewCreate(http.DefaultClient)
	cmd := &cli.Command{
		Commands: []*cli.Command{
			createCmd.Command(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("error running command", "error", err)
		os.Exit(1)
	}
}
