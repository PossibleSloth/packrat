package main

import (
	"context"
	"fmt"
	"os"

	"github.com/possiblesloth/packrat/app"
)

func main() {
	ctx := context.Background()
	if err := app.Run(ctx, os.Stdout, os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
