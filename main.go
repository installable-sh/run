package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/installable-sh/run/internal/run"
	"mvdan.cc/sh/v3/interp"
)

func main() {
	// Set up signal handling - catch SIGINT/SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	var receivedSig os.Signal
	go func() {
		receivedSig = <-sigCh
		cancel()
	}()

	cmd := run.New(os.Args[1:])
	if err := cmd.Exec(ctx); err != nil {
		// Check if shell returned an exit status
		var exitStatus interp.ExitStatus
		if errors.As(err, &exitStatus) {
			os.Exit(int(exitStatus))
		}

		// If we received a signal, exit with 128 + signal number
		if receivedSig != nil {
			if sig, ok := receivedSig.(syscall.Signal); ok {
				os.Exit(128 + int(sig))
			}
		}

		// Don't print error if we were interrupted by signal
		if ctx.Err() == nil {
			fmt.Fprintf(os.Stderr, "[run] error: %v\n", err)
		}
		os.Exit(1)
	}
}
