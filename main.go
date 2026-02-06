package main

import (
	"context"
	"fmt"
	"os"

	"github.com/installable-sh/run/internal/run"
)

func main() {
	ctx := context.Background()

	cmd := run.New(os.Args[1:])
	if err := cmd.Exec(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "[run] error: %v\n", err)
		os.Exit(1)
	}
}
