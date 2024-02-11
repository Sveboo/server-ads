package tests

import (
	"ads-server/internal/adapters/repo"
	"ads-server/internal/app"
	grpcPort "ads-server/internal/ports/grpc"
	grpc2 "ads-server/proto"
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
	"time"
)

// Table tests

type table struct {
	id          int64
	name, email string
}

func TestCreateUsers(t *testing.T) {

	var testTable = [...]table{
		{0, "Jonh", "os"},
		{1, "Mary", "tos"},
		{2, "Craig", "nos"},
		{3, "Kate", "ros"},
	}

	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewAdService(app.NewApp(repo.NewAd(), repo.NewUser()))
	grpc2.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpc2.NewAdServiceClient(conn)

	for _, test := range testTable {
		user, err := client.CreateUser(ctx, &grpc2.CreateUserRequest{
			Name:  test.name,
			Email: test.email,
		})
		assert.NoError(t, err)
		assert.Equal(t, test.id, user.Id)
		assert.Equal(t, test.name, user.Name)
		assert.Equal(t, test.email, user.Email)
	}
}
