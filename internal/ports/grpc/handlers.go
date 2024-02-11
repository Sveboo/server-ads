package grpc

import (
	"ads-server/internal/app"
	"ads-server/internal/errs"
	proto "ads-server/proto"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name IAdService
type IAdService interface {
	CreateAd(ctx context.Context, request *proto.CreateAdRequest) (*proto.AdResponse, error)
	ChangeAdStatus(ctx context.Context, request *proto.ChangeAdStatusRequest) (*proto.AdResponse, error)
	UpdateAd(ctx context.Context, request *proto.UpdateAdRequest) (*proto.AdResponse, error)
	ListAds(ctx context.Context, request *proto.ListAdRequest) (*proto.ListAdResponse, error)
	CreateUser(ctx context.Context, request *proto.CreateUserRequest) (*proto.UserResponse, error)
	GetUser(ctx context.Context, request *proto.GetUserRequest) (*proto.UserResponse, error)
	UpdateUser(ctx context.Context, request *proto.UpdateUserRequest) (*proto.UserResponse, error)
	DeleteAd(ctx context.Context, request *proto.DeleteAdRequest) (*proto.DeleteAdResponse, error)
}
type AdService struct {
	app app.IApp
}

func NewAdService(a app.App) *AdService {
	return &AdService{app: a}
}

func (a *AdService) CreateAd(ctx context.Context, request *proto.CreateAdRequest) (*proto.AdResponse, error) {
	if _, err := a.app.FindUser(ctx, request.UserId); errors.Is(err, errs.UserNotFoundError) {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	ad, err := a.app.CreateAd(ctx, request.UserId, request.Title, request.Text)
	if errors.Is(err, errs.ValidationError) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if errors.Is(err, errs.AccessError) {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	return &proto.AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		AuthorId:  ad.AuthorID,
		Published: ad.Published,
	}, nil
}

func (a *AdService) ChangeAdStatus(ctx context.Context, request *proto.ChangeAdStatusRequest) (*proto.AdResponse, error) {

	if _, err := a.app.FindUser(ctx, request.UserId); errors.Is(err, errs.UserNotFoundError) {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	ad, err := a.app.PublishAd(ctx, request.AdId, request.UserId, request.Published)

	if errors.Is(err, errs.AccessError) {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &proto.AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		AuthorId:  ad.AuthorID,
		Published: ad.Published,
	}, nil
}

func (a *AdService) UpdateAd(ctx context.Context, request *proto.UpdateAdRequest) (*proto.AdResponse, error) {

	if _, err := a.app.FindUser(ctx, request.UserId); errors.Is(err, errs.UserNotFoundError) {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	ad, err := a.app.UpdateAd(ctx, request.AdId, request.UserId, request.Title, request.Text)

	if errors.Is(err, errs.ValidationError) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if errors.Is(err, errs.AccessError) {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	return &proto.AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		AuthorId:  ad.AuthorID,
		Published: ad.Published,
	}, nil
}

func (a *AdService) ListAds(ctx context.Context, request *proto.ListAdRequest) (*proto.ListAdResponse, error) {
	ads := a.app.GetAdByName(ctx, request.Title)
	if len(ads) == 0 {
		return nil, status.Error(codes.NotFound, errs.AdNotFoundError.Error())
	}

	list := make([]*proto.AdResponse, len(ads))
	for i, ad := range ads {
		list[i] = &proto.AdResponse{
			Id:        ad.ID,
			Title:     ad.Title,
			Text:      ad.Text,
			AuthorId:  ad.AuthorID,
			Published: ad.Published,
		}
	}

	return &proto.ListAdResponse{
		List: list,
	}, nil
}

func (a *AdService) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (*proto.UserResponse, error) {
	user, err := a.app.CreateUser(ctx, request.Name, request.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &proto.UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (a *AdService) GetUser(ctx context.Context, request *proto.GetUserRequest) (*proto.UserResponse, error) {
	if request.Id == nil {
		return nil, errs.WrongProtoBufDataError
	}
	user, err := a.app.FindUser(ctx, request.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &proto.UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (a *AdService) DeleteUser(ctx context.Context, request *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {

	err := a.app.DeleteUser(ctx, request.Id)
	if err != nil {
		return &proto.DeleteUserResponse{Success: false}, status.Error(codes.NotFound, err.Error())
	}

	return &proto.DeleteUserResponse{Success: true}, nil

}

func (a *AdService) UpdateUser(ctx context.Context, request *proto.UpdateUserRequest) (*proto.UserResponse, error) {
	user, err := a.app.UpdateUser(ctx, request.Id, request.Name, request.Email)
	if err != nil {
		return nil, err
	}

	return &proto.UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (a *AdService) DeleteAd(ctx context.Context, request *proto.DeleteAdRequest) (*proto.DeleteAdResponse, error) {

	err := a.app.DeleteAd(ctx, request.AdId, request.AuthorId)
	if errors.Is(err, errs.AdNotFoundError) {
		return &proto.DeleteAdResponse{Success: false}, status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, errs.AccessError) {
		return &proto.DeleteAdResponse{Success: false}, status.Error(codes.PermissionDenied, err.Error())
	}

	return &proto.DeleteAdResponse{Success: true}, nil
}
