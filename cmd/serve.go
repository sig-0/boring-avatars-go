package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/sig-0/boring-avatars-go/server"
	"github.com/sig-0/boring-avatars-go/server/config"
	"golang.org/x/sync/errgroup"
)

const envPrefix = "BORING_AVATARS"

// serveCfg wraps the serve configuration
type serveCfg struct {
	config *config.Config

	configPath string
}

// newServeCmd creates the serve command
func newServeCmd() *ffcli.Command {
	cfg := &serveCfg{
		config: config.DefaultConfig(),
	}

	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	cfg.registerFlags(fs)

	return &ffcli.Command{
		Name:       "serve",
		ShortUsage: "serve [flags]",
		LongHelp:   "Serves the Boring Avatars server",
		FlagSet:    fs,
		Exec:       cfg.exec,
		Options: []ff.Option{
			// Allow using ENV variables
			ff.WithEnvVars(),
			ff.WithEnvVarPrefix(envPrefix),
		},
	}
}

// registerFlags registers the serve command flags
func (c *serveCfg) registerFlags(fs *flag.FlagSet) {
	fs.StringVar(
		&c.config.ListenAddress,
		"listen",
		config.DefaultListenAddress,
		"the IP:PORT URL for the server",
	)

	fs.StringVar(
		&c.configPath,
		"config",
		"",
		"the path to the server TOML configuration, if any",
	)
}

// exec executes the server serve command
func (c *serveCfg) exec(ctx context.Context, _ []string) error {
	// Read the server configuration, if any
	if c.configPath != "" {
		serverCfg, err := config.Read(c.configPath)
		if err != nil {
			return fmt.Errorf("unable to read server config, %w", err)
		}

		c.config = serverCfg
	}

	// Create a new logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create the server instance
	s, err := server.New(
		server.WithLogger(logger),
		server.WithConfig(c.config),
	)
	if err != nil {
		return fmt.Errorf("unable to create server, %w", err)
	}

	runCtx, cancelFn := signal.NotifyContext(
		ctx,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	defer cancelFn()

	group, gCtx := errgroup.WithContext(runCtx)

	// Start the HTTP server
	group.Go(func() error {
		return s.Serve(gCtx)
	})

	return group.Wait()
}
