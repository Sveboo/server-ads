package tests

import (
	grpc2 "ads-server/proto"
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"

	"ads-server/internal/adapters/repo"
	"ads-server/internal/app"
	grpcPort "ads-server/internal/ports/grpc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestGRPCCreateUser(t *testing.T) {
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
	res, err := client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Name)
}

func TestGRPCUpdateUser(t *testing.T) {
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
	_, err = client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	res, err := client.UpdateUser(ctx, &grpc2.UpdateUserRequest{
		Id:    0,
		Name:  "Vanya",
		Email: "mmm@mail.go",
	})

	assert.NoError(t, err)
	assert.Equal(t, "Vanya", res.Name)
	assert.Equal(t, "mmm@mail.go", res.Email)
}

func TestGRPCGetExistingUser(t *testing.T) {
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

	user, err := client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Oleg"})
	if err != nil {
		return
	}

	id := int64(0)
	res, err := client.GetUser(ctx, &grpc2.GetUserRequest{Id: &id})
	assert.NoError(t, err, "client.GetUser")
	assert.Equal(t, user.Id, res.Id)
	assert.Equal(t, user.Name, res.Name)
	assert.Equal(t, user.Email, res.Email)
}

func TestGRPCGetNonexistentUser(t *testing.T) {
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

	id := int64(0)
	res, err := client.GetUser(ctx, &grpc2.GetUserRequest{Id: &id})
	assert.Emptyf(t, res, "no user should be found")
	assert.Error(t, err, "error during GetUser expected")
}

func TestGRPCDeleteUser(t *testing.T) {
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
	_, err = client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Oleg"})
	if err != nil {
		return
	}

	res, err := client.DeleteUser(ctx, &grpc2.DeleteUserRequest{Id: 0})
	assert.NoError(t, err, "client.DeleteUser")
	assert.Equal(t, true, res.Success)

	var id int64
	user, err := client.GetUser(ctx, &grpc2.GetUserRequest{Id: &id})
	assert.Error(t, err)
	assert.Empty(t, user)
}

func TestGRPCCreateAd(t *testing.T) {
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

	_, err = client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Oleg"})
	if err != nil {
		return
	}

	res, err := client.CreateAd(ctx, &grpc2.CreateAdRequest{
		Title:  "Hello",
		Text:   "World",
		UserId: 0,
	})

	assert.Equal(t, "Hello", res.Title)
	assert.Equal(t, "World", res.Text)
	assert.Equal(t, int64(0), res.AuthorId)
	assert.NoError(t, err)
}

func TestGRPCChangeAdStatus(t *testing.T) {
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

	_, err = client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Oleg"})
	if err != nil {
		return
	}

	_, err = client.CreateAd(ctx, &grpc2.CreateAdRequest{
		Title:  "Hello",
		Text:   "World",
		UserId: 0,
	})

	if err != nil {
		return
	}

	res, err := client.ChangeAdStatus(ctx, &grpc2.ChangeAdStatusRequest{
		AdId:      0,
		UserId:    0,
		Published: true,
	})

	assert.Equal(t, true, res.Published)
	assert.NoError(t, err)
}

func TestGRPCChangeAdStatusFromAnotherUser(t *testing.T) {
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

	_, err = client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Oleg"})
	if err != nil {
		return
	}

	_, err = client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Anton"})
	if err != nil {
		return
	}

	_, err = client.CreateAd(ctx, &grpc2.CreateAdRequest{
		Title:  "Hello",
		Text:   "World",
		UserId: 0,
	})
	if err != nil {
		return
	}

	res, err := client.ChangeAdStatus(ctx, &grpc2.ChangeAdStatusRequest{
		AdId:      1,
		UserId:    0,
		Published: true,
	})

	assert.Empty(t, res)
	assert.Error(t, err)
}

func TestGRPCUpdateAd(t *testing.T) {
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

	_, err = client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Oleg"})
	if err != nil {
		return
	}

	_, err = client.CreateAd(ctx, &grpc2.CreateAdRequest{
		Title:  "Hello",
		Text:   "World",
		UserId: 0,
	})
	if err != nil {
		return
	}

	res, err := client.UpdateAd(ctx, &grpc2.UpdateAdRequest{
		AdId:   0,
		Title:  "I'm",
		Text:   "Updated",
		UserId: 0,
	})

	assert.Equal(t, "I'm", res.Title)
	assert.Equal(t, "Updated", res.Text)
	assert.NoError(t, err)
}

func TestGRPCDeleteAd(t *testing.T) {
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

	_, err = client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Oleg"})
	if err != nil {
		return
	}

	_, err = client.CreateAd(ctx, &grpc2.CreateAdRequest{
		Title:  "Hello",
		Text:   "World",
		UserId: 0,
	})
	if err != nil {
		return
	}

	res, err := client.DeleteAd(ctx, &grpc2.DeleteAdRequest{
		AdId:     0,
		AuthorId: 0,
	})

	assert.NoError(t, err)
	assert.Equal(t, true, res.Success)

	ads, err := client.ListAds(ctx, &grpc2.ListAdRequest{Title: ""})
	assert.Error(t, err)
	assert.Empty(t, ads)
}

func TestGRPCListAdsNoTitle(t *testing.T) {
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

	_, err = client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Oleg"})
	if err != nil {
		return
	}

	_, err = client.CreateAd(ctx, &grpc2.CreateAdRequest{
		Title:  "Hello",
		Text:   "World",
		UserId: 0,
	})
	if err != nil {
		return
	}

	_, err = client.CreateAd(ctx, &grpc2.CreateAdRequest{
		Title:  "Hello1",
		Text:   "World1",
		UserId: 0,
	})
	if err != nil {
		return
	}
	_, err = client.CreateAd(ctx, &grpc2.CreateAdRequest{
		Title:  "Hello2",
		Text:   "World2",
		UserId: 0,
	})
	if err != nil {
		return
	}

	_, err = client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Misha"})
	if err != nil {
		return
	}

	_, err = client.CreateAd(ctx, &grpc2.CreateAdRequest{
		Title:  "Hello3",
		Text:   "World3",
		UserId: 1,
	})
	if err != nil {
		return
	}

	res, err := client.ListAds(ctx, &grpc2.ListAdRequest{Title: ""})

	assert.Error(t, err)
	assert.Empty(t, res)

	_, err = client.ChangeAdStatus(ctx, &grpc2.ChangeAdStatusRequest{
		AdId:      0,
		UserId:    0,
		Published: true,
	})
	if err != nil {
		return
	}

	_, err = client.ChangeAdStatus(ctx, &grpc2.ChangeAdStatusRequest{
		AdId:      1,
		UserId:    0,
		Published: true,
	})
	if err != nil {
		return
	}

	res, err = client.ListAds(ctx, &grpc2.ListAdRequest{Title: ""})
	assert.Len(t, res.List, 2)
	assert.NoError(t, err)
}

func TestGRPCListAdsWithTitle(t *testing.T) {
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

	_, err = client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Oleg"})
	if err != nil {
		return
	}

	_, err = client.CreateAd(ctx, &grpc2.CreateAdRequest{
		Title:  "Hello",
		Text:   "World",
		UserId: 0,
	})
	if err != nil {
		return
	}

	_, err = client.CreateAd(ctx, &grpc2.CreateAdRequest{
		Title:  "Bye",
		Text:   "World1",
		UserId: 0,
	})
	if err != nil {
		return
	}
	_, err = client.CreateAd(ctx, &grpc2.CreateAdRequest{
		Title:  "Hello2",
		Text:   "World2",
		UserId: 0,
	})
	if err != nil {
		return
	}

	_, err = client.CreateUser(ctx, &grpc2.CreateUserRequest{Name: "Misha"})
	if err != nil {
		return
	}

	_, err = client.CreateAd(ctx, &grpc2.CreateAdRequest{
		Title:  "Hello3",
		Text:   "World3",
		UserId: 1,
	})
	if err != nil {
		return
	}

	res, err := client.ListAds(ctx, &grpc2.ListAdRequest{Title: ""})

	assert.Error(t, err)
	assert.Empty(t, res)

	_, err = client.ChangeAdStatus(ctx, &grpc2.ChangeAdStatusRequest{
		AdId:      0,
		UserId:    0,
		Published: true,
	})
	if err != nil {
		return
	}

	_, err = client.ChangeAdStatus(ctx, &grpc2.ChangeAdStatusRequest{
		AdId:      1,
		UserId:    0,
		Published: true,
	})
	if err != nil {
		return
	}

	res, err = client.ListAds(ctx, &grpc2.ListAdRequest{Title: "Hello"})
	assert.Len(t, res.List, 1)
	assert.NoError(t, err)
}
