package grpc

import (
	"ads-server/internal/adapters/repo"
	"ads-server/internal/ports/grpc/pkg/interceptors"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"reflect"
	"testing"
)

func TestNewGRPCServer(t *testing.T) {
	adRepo := repo.NewAd()
	userRepo := repo.NewUser()
	recoveryOpt := []grpcrecovery.Option{
		grpcrecovery.WithRecoveryHandler(interceptors.RecoveryFunc),
	}
	want := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors.Logger,
		grpcrecovery.UnaryServerInterceptor(recoveryOpt...)))

	t.Run("test correct output", func(t *testing.T) {
		if got := NewGRPCServer(adRepo, userRepo); !(reflect.TypeOf(want).Kind() == reflect.TypeOf(got).Kind()) {
			t.Errorf("NewGRPCServer() = %v, want %v", got, want)
		}
	})

}
