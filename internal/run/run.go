package run

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/installable-sh/lib/fetch"
	"github.com/installable-sh/lib/shell"
	"github.com/installable-sh/lib/version"
)

// Run represents the RUN command with parsed arguments.
type Run struct {
	ShowHelp    bool
	ShowVersion bool
	SendEnv     bool
	Raw         bool
	NoCache     bool
	URL         string
	ScriptArgs  []string

	// IO streams (defaults to os.Stdin/Stdout/Stderr)
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// New parses command-line arguments and returns a Run command.
func New(args []string) *Run {
	r := &Run{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	// Find the URL (first arg starting with http:// or https://)
	urlIndex := -1
	for i, arg := range args {
		if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
			urlIndex = i
			break
		}
	}

	var runArgs []string
	if urlIndex >= 0 {
		runArgs = args[:urlIndex]
		r.URL = args[urlIndex]
		r.ScriptArgs = args[urlIndex+1:]
	} else {
		runArgs = args
	}

	for _, arg := range runArgs {
		switch arg {
		case "--help", "-h":
			r.ShowHelp = true
		case "--version", "-v":
			r.ShowVersion = true
		case "+env":
			r.SendEnv = true
		case "+raw":
			r.Raw = true
		case "+nocache":
			r.NoCache = true
		}
	}

	return r
}

// Exec executes the RUN command.
func (r *Run) Exec(ctx context.Context) error {
	if r.ShowVersion {
		version.Print("RUN")
		return nil
	}

	if r.ShowHelp || r.URL == "" {
		_, _ = fmt.Fprintln(r.Stdout, "usage: RUN [+env] [+raw] [+nocache] <url> [args...]")
		_, _ = fmt.Fprintln(r.Stdout, "  +env      Send environment variables as X-Env-* headers")
		_, _ = fmt.Fprintln(r.Stdout, "  +raw      Print the script without executing")
		_, _ = fmt.Fprintln(r.Stdout, "  +nocache  Bypass CDN caches")
		return nil
	}

	client, err := fetch.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	script, err := fetch.Fetch(ctx, client, fetch.Options{
		URL:     r.URL,
		SendEnv: r.SendEnv,
		NoCache: r.NoCache,
	})
	if err != nil {
		return fmt.Errorf("failed to fetch script: %w", err)
	}

	if r.Raw {
		_, _ = fmt.Fprint(r.Stdout, script.Content)
		return nil
	}

	return shell.RunWithIO(
		ctx,
		shell.Script{Content: script.Content, Name: script.Name},
		r.ScriptArgs,
		r.Stdin,
		r.Stdout,
		r.Stderr,
	)
}
