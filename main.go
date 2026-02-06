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
	err := cmd.Exec(ctx)

	// No error - clean exit
	if err == nil {
		os.Exit(0)
	}

	// Check if shell returned an exit status
	var exitStatus interp.ExitStatus
	if errors.As(err, &exitStatus) {
		os.Exit(int(exitStatus))
	}

	// If context was cancelled due to signal, exit with 128 + signal number
	if errors.Is(err, context.Canceled) && receivedSig != nil {
		if sig, ok := receivedSig.(syscall.Signal); ok {
			os.Exit(128 + int(sig))
		}
	}

	// Other error
	fmt.Fprintf(os.Stderr, "[run] error: %v\n", err)
	os.Exit(1)
}
