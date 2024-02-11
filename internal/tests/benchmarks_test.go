package tests

import (
	"ads-server/internal/adapters/repo"
	"ads-server/internal/app"
	grpcPort "ads-server/internal/ports/grpc"
	grpc2 "ads-server/proto"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
	"time"
)

func BenchmarkGRPCRequest(b *testing.B) {
	lis := bufconn.Listen(1024 * 1024)
	b.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	b.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewAdService(app.NewApp(repo.NewAd(), repo.NewUser()))
	grpc2.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(b, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	b.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(b, err, "grpc.DialContext")

	b.Cleanup(func() {
		conn.Close()
	})

	client := grpc2.NewAdServiceClient(conn)

	_, err = client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Oleg"})
	if err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.CreateAd(ctx, &grpc2.CreateAdRequest{
			Title:  "Hello",
			Text:   "World",
			UserId: int64(i),
		})
	}
}

func BenchmarkHTTPRequest(b *testing.B) {
	client := getTestClient()
	_, err := client.createUser(0, "James", "Ostin")
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < b.N; i++ {
		_, _ = client.createAd(int64(i), "hello", "world")
	}
}
