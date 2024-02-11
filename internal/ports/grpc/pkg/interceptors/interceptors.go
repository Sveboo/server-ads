package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()
	i, err := handler(ctx, req)
	end := time.Since(start).Nanoseconds()

	if err != nil {
		log.Printf("Error during request to %s occured: %v\nRequest aborted", info.FullMethod, err)
		return nil, err
	}

	log.Printf("Successful request to %s\nElapsed time: %d ns", info.FullMethod, end)
	return i, nil
}

func RecoveryFunc(p any) error {
	return status.Errorf(codes.Unknown, "panic triggered: %v", p)
}
