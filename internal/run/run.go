package run

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/installable-sh/lib/fetch"
	"github.com/installable-sh/lib/log"
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
	Debug       bool
	URL         string
	ScriptArgs  []string

	// IO streams (defaults to os.Stdin/Stdout/Stderr)
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// New parses command-line arguments and returns a Run command.
// Flags can also be set via environment variables:
// PLUS_ENV=true, PLUS_RAW=true, PLUS_NOCACHE=true, PLUS_DEBUG=true
func New(args []string) *Run {
	r := &Run{
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		SendEnv: envBool("PLUS_ENV"),
		Raw:     envBool("PLUS_RAW"),
		NoCache: envBool("PLUS_NOCACHE"),
		Debug:   envBool("PLUS_DEBUG"),
	}

	// Find the URL (first arg starting with http://, https://, or file://)
	urlIndex := -1
	for i, arg := range args {
		if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") || strings.HasPrefix(arg, "file://") {
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
		case "+debug":
			r.Debug = true
		}
	}

	return r
}

// envBool returns true if the environment variable is set to "true", "1", or "yes".
func envBool(name string) bool {
	v := strings.ToLower(os.Getenv(name))
	return v == "true" || v == "1" || v == "yes"
}

// Exec executes the RUN command.
func (r *Run) Exec(ctx context.Context) error {
	if r.ShowVersion {
		version.Print("RUN")
		return nil
	}

	if r.ShowHelp || r.URL == "" {
		_, _ = fmt.Fprintln(r.Stdout, "usage: RUN [+env] [+raw] [+nocache] [+debug] <url> [args...]")
		_, _ = fmt.Fprintln(r.Stdout, "  +env      Send environment variables as X-Env-* headers")
		_, _ = fmt.Fprintln(r.Stdout, "  +raw      Print the script without executing")
		_, _ = fmt.Fprintln(r.Stdout, "  +nocache  Bypass CDN caches")
		_, _ = fmt.Fprintln(r.Stdout, "  +debug    Show debug information (URL, headers, response)")
		return nil
	}

	logger := log.New("run")
	logger.SetDebug(r.Debug)
	logger.SetOutput(r.Stderr)

	var script fetch.Script

	logger.Debugf("URL: %s, Flags: +env=%v +raw=%v +nocache=%v", r.URL, r.SendEnv, r.Raw, r.NoCache)

	if filePath, ok := strings.CutPrefix(r.URL, "file://"); ok {
		// Handle file:// URLs by reading from local filesystem
		logger.Debugf("Reading local file: %s", filePath)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return logger.Errorf("failed to read %s: %w", r.URL, err)
		}
		script = fetch.Script{
			Content: string(content),
			Name:    filepath.Base(filePath),
		}
	} else {
		// Handle http:// and https:// URLs
		client, err := fetch.NewClient(logger)
		if err != nil {
			return logger.Errorf("failed to create HTTP client: %w", err)
		}

		script, err = fetch.Fetch(ctx, client, fetch.Options{
			URL:     r.URL,
			SendEnv: r.SendEnv,
			NoCache: r.NoCache,
		}, logger)
		if err != nil {
			// Error already logged by fetch
			return err
		}
	}

	if r.Raw {
		_, _ = fmt.Fprint(r.Stdout, script.Content)
		return nil
	}

	return shell.Run(
		ctx,
		shell.Script{Content: script.Content, Name: script.Name},
		r.ScriptArgs,
		r.Stdin,
		r.Stdout,
		r.Stderr,
		logger,
	)
}
