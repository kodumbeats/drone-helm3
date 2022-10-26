package run

import (
	"fmt"
	"github.com/mongodb-forks/drone-helm3/internal/env"
)

// List is a step in a helm Plan that calls `helm help`.
type List struct {
	*config
	helmCommand string
	cmd         cmd
}

// NewList creates a List using fields from the given Config. No validation is performed at this time.
func NewList(cfg env.Config) *List {
	return &List{
		config:      newConfig(cfg),
		helmCommand: cfg.Command,
	}
}

// Execute executes the `helm help` command.
func (l *List) Execute() error {
	if err := l.cmd.Run(); err != nil {
		return fmt.Errorf("while running '%s': %w", l.cmd.String(), err)
	}

	if l.helmCommand == "list" {
		return nil
	}
	return fmt.Errorf("unknown command '%s'", l.helmCommand)
}

// Prepare gets the List ready to execute.
func (l *List) Prepare() error {
	args := l.globalFlags()
	args = append(args, "list", "--output", "json")

	l.cmd = command(helmBin, args...)
	l.cmd.Stdout(l.stdout)
	l.cmd.Stderr(l.stderr)

	if l.debug {
		fmt.Fprintf(l.stderr, "Generated command: '%s'\n", l.cmd.String())
	}

	return nil
}
