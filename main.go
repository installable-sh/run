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

	exitCode := 0
	cmd := run.New(os.Args[1:])

	if err := cmd.Exec(ctx); err != nil {
		exitCode = 1

		// Check if shell returned an exit status
		var exitStatus interp.ExitStatus
		if errors.As(err, &exitStatus) {
			exitCode = int(exitStatus)
		} else if errors.Is(err, context.Canceled) && receivedSig != nil {
			// Context cancelled due to signal - use 128 + signal number
			if sig, ok := receivedSig.(syscall.Signal); ok {
				exitCode = 128 + int(sig)
			}
		} else {
			fmt.Fprintf(os.Stderr, "[run] error: %v\n", err)
		}
	}

	os.Exit(exitCode)
}
