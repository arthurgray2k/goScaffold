package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	slog.Info("Worker started. Press Ctrl+C to stop.")

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				slog.Info("Processing background job...")
			}
		}
	}()

	<-quit
	slog.Info("Shutting down worker gracefully...")
	cancel()
	
	// Give worker time to finish current job
	time.Sleep(1 * time.Second)
	slog.Info("Worker stopped.")
}
