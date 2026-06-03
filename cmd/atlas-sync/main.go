package main

import (
	"log/slog"
	"os"

	"github.com/RomanUrakhov/atlas-sync/cmd/atlas-sync/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error("sync failed", "err", err)
		os.Exit(1)
	}
}
