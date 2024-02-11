package grpc

import (
	"ads-server/internal/app"
	"ads-server/internal/ports/grpc/pkg/interceptors"
	proto "ads-server/proto"
	"context"
	"fmt"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewGRPCServer(a app.AdRepository, u app.UserRepository) *grpc.Server {

	service := &AdService{
		app: app.NewApp(a, u),
	}

	recoveryOpt := []grpcrecovery.Option{
		grpcrecovery.WithRecoveryHandler(interceptors.RecoveryFunc),
	}

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors.Logger,
		grpcrecovery.UnaryServerInterceptor(recoveryOpt...)))
	proto.RegisterAdServiceServer(server, service)

	return server
}

// Run returns function to start gRPC server on a port given and implements graceful shutdown principle
func Run(ctx context.Context, a app.AdRepository, u app.UserRepository, grpcPort string) func() error {
	return func() error {
		grpcServer := NewGRPCServer(a, u)

		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		errCh := make(chan error)

		defer func() {
			grpcServer.GracefulStop()
			_ = lis.Close()

			close(errCh)
		}()

		go func() {
			if err = grpcServer.Serve(lis); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err = <-errCh:
			return fmt.Errorf("grpc server can't listen and serve requests: %w", err)
		}
	}
}
