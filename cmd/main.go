package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	fs := flag.NewFlagSet("root", flag.ExitOnError)

	// Create the root command
	cmd := &ffcli.Command{
		ShortUsage: "<sub-command> [flags] [<arg>...]",
		LongHelp:   "Runs the boring avatar generation service",
		FlagSet:    fs,
		Exec: func(_ context.Context, _ []string) error {
			return flag.ErrHelp
		},
	}

	// Add the subcommands
	cmd.Subcommands = []*ffcli.Command{
		newServeCmd(),
		newGenerateCmd(),
	}

	if err := cmd.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
