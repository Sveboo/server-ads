package main

import (
	"ads-server/internal/adapters/repo"
	"ads-server/internal/ports/grpc"
	"ads-server/internal/ports/httpgin"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	grpcPort = ":50054"
	httpPort = ":8080"
)

func captureSigQuit(ctx context.Context) func() error {
	return func() error {
		sigQuit := make(chan os.Signal, 1)
		signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
		signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v ", s)
		case <-ctx.Done():
			return nil
		}
	}
}

func main() {
	a := repo.NewAd()
	u := repo.NewUser()
	eg, ctx := errgroup.WithContext(context.Background())

	// capture signals to stop working
	eg.Go(captureSigQuit(ctx))

	// run gRPC server
	eg.Go(grpc.Run(ctx, a, u, grpcPort))

	// run HTTP server
	eg.Go(httpgin.Run(ctx, a, u, httpPort))

	err := eg.Wait()
	if err != nil {
		fmt.Println(err)
		return
	}
}
