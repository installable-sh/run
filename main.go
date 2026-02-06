package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/installable-sh/run/internal/run"
)

func main() {
	// Set up signal handling - catch SIGINT/SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	cmd := run.New(os.Args[1:])
	if err := cmd.Exec(ctx); err != nil {
		// Don't print error if we were interrupted by signal
		if ctx.Err() == nil {
			fmt.Fprintf(os.Stderr, "[run] error: %v\n", err)
		}
		os.Exit(1)
	}
}
