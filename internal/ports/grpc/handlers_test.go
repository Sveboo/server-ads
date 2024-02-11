package grpc

import (
	"ads-server/internal/adapters/repo"
	"ads-server/internal/ads"
	"ads-server/internal/app"
	"ads-server/internal/errs"
	"ads-server/internal/users"
	"ads-server/mocks"
	proto "ads-server/proto"
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestAdService_ChangeAdStatus(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *proto.ChangeAdStatusRequest
	}
	tests := []struct {
		name      string
		args      args
		want      *proto.AdResponse
		wantErr   bool
		userExist error
		adError   error
	}{
		{
			name: "No Error",
			args: args{
				ctx: context.Background(),
				request: &proto.ChangeAdStatusRequest{
					AdId:      0,
					UserId:    0,
					Published: true,
				},
			},
			want: &proto.AdResponse{
				Id:        0,
				Title:     "example",
				Text:      "example",
				AuthorId:  0,
				Published: true,
			},
			wantErr:   false,
			userExist: nil,
			adError:   nil,
		},

		{
			name: "No such user",
			args: args{
				ctx: context.Background(),
				request: &proto.ChangeAdStatusRequest{
					AdId:      0,
					UserId:    1,
					Published: true,
				},
			},
			want:      nil,
			wantErr:   true,
			userExist: errs.UserNotFoundError,
			adError:   nil,
		},

		{
			name: "Validation error",
			args: args{
				ctx: context.Background(),
				request: &proto.ChangeAdStatusRequest{
					AdId:      0,
					UserId:    1,
					Published: true,
				},
			},
			want:      nil,
			wantErr:   true,
			userExist: nil,
			adError:   errs.AccessError,
		},

		{
			name: "Unknown error",
			args: args{
				ctx: context.Background(),
				request: &proto.ChangeAdStatusRequest{
					AdId:      0,
					UserId:    1,
					Published: true,
				},
			},
			want:      nil,
			wantErr:   true,
			userExist: nil,
			adError:   errors.New("some error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeApp := mocks.NewIApp(t)
			fakeApp.
				On("FindUser", tt.args.ctx, tt.args.request.UserId).
				Return(nil, tt.userExist).
				Maybe()
			fakeApp.
				On("PublishAd", tt.args.ctx, tt.args.request.AdId, tt.args.request.UserId,
					tt.args.request.Published).
				Return(&ads.Ad{
					ID:        0,
					Title:     "example",
					Text:      "example",
					AuthorID:  tt.args.request.UserId,
					CDate:     time.Time{},
					UDate:     time.Time{},
					Published: tt.args.request.Published,
				}, tt.adError).
				Maybe()
			a := &AdService{
				app: fakeApp,
			}
			got, err := a.ChangeAdStatus(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeAdStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChangeAdStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdService_CreateAd(t *testing.T) {

	type args struct {
		ctx     context.Context
		request *proto.CreateAdRequest
	}
	tests := []struct {
		name      string
		args      args
		want      *proto.AdResponse
		userExist error
		adExist   error
		wantErr   bool
	}{
		{
			name: "No Error",
			args: args{
				ctx: context.Background(),
				request: &proto.CreateAdRequest{
					Title:  "hello",
					Text:   "World",
					UserId: 0,
				},
			},
			want: &proto.AdResponse{
				Id:        0,
				Title:     "hello",
				Text:      "World",
				AuthorId:  0,
				Published: false,
			},
			wantErr:   false,
			userExist: nil,
			adExist:   nil,
		},

		{
			name: "Error: no such user",
			args: args{
				ctx: context.Background(),
				request: &proto.CreateAdRequest{
					Title:  "hello",
					Text:   "World",
					UserId: 1,
				},
			},
			want:      nil,
			wantErr:   true,
			userExist: errs.UserNotFoundError,
			adExist:   nil,
		},

		{
			name: "Error: validation",
			args: args{
				ctx: context.Background(),
				request: &proto.CreateAdRequest{
					Title:  "",
					Text:   "",
					UserId: 0,
				},
			},
			want:      nil,
			wantErr:   true,
			userExist: nil,
			adExist:   errs.ValidationError,
		},

		{
			name: "Error: storage overflow",
			args: args{
				ctx: context.Background(),
				request: &proto.CreateAdRequest{
					Title:  "hello",
					Text:   "World",
					UserId: 0,
				},
			},
			want:      nil,
			wantErr:   true,
			userExist: nil,
			adExist:   errs.AccessError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeApp := mocks.NewIApp(t)
			fakeApp.
				On("FindUser", tt.args.ctx, tt.args.request.UserId).
				Return(nil, tt.userExist).
				Maybe()
			fakeApp.
				On("CreateAd", tt.args.ctx, tt.args.request.UserId,
					tt.args.request.Title, tt.args.request.Text).
				Return(&ads.Ad{
					ID:        0,
					Title:     tt.args.request.Title,
					Text:      tt.args.request.Text,
					AuthorID:  0,
					CDate:     time.Time{},
					UDate:     time.Time{},
					Published: false,
				}, tt.adExist).
				Maybe()
			a := &AdService{
				app: fakeApp,
			}
			got, err := a.CreateAd(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateAd() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdService_CreateUser(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *proto.CreateUserRequest
	}
	tests := []struct {
		name      string
		args      args
		want      *proto.UserResponse
		userExist error
		wantErr   bool
	}{
		{
			name: "No Error",
			args: args{
				ctx: context.Background(),
				request: &proto.CreateUserRequest{
					Name:  "James",
					Email: "Ostin",
				},
			},
			want: &proto.UserResponse{
				Id:    0,
				Name:  "James",
				Email: "Ostin",
			},
			wantErr:   false,
			userExist: nil,
		},

		{
			name: "overflow",
			args: args{
				ctx: context.Background(),
				request: &proto.CreateUserRequest{
					Name:  "Jeremy",
					Email: "Hopkins",
				},
			},
			want:      nil,
			wantErr:   true,
			userExist: errs.UserNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeApp := mocks.NewIApp(t)
			fakeApp.
				On("CreateUser", tt.args.ctx, tt.args.request.Name, tt.args.request.Email).
				Return(&users.User{
					ID:    0,
					Name:  tt.args.request.Name,
					Email: tt.args.request.Email,
				}, tt.userExist).
				Maybe()
			a := &AdService{
				app: fakeApp,
			}
			got, err := a.CreateUser(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdService_DeleteAd(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *proto.DeleteAdRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *proto.DeleteAdResponse
		adExist error
		wantErr bool
	}{
		{
			name: "No Error",
			args: args{
				ctx: context.Background(),
				request: &proto.DeleteAdRequest{
					AdId:     0,
					AuthorId: 0,
				},
			},
			want:    &proto.DeleteAdResponse{Success: true},
			wantErr: false,
			adExist: nil,
		},

		{
			name: "Access error",
			args: args{
				ctx: context.Background(),
				request: &proto.DeleteAdRequest{
					AdId:     0,
					AuthorId: 1,
				},
			},
			want:    &proto.DeleteAdResponse{Success: false},
			wantErr: true,
			adExist: errs.AccessError,
		},

		{
			name: "no such ad",
			args: args{
				ctx: context.Background(),
				request: &proto.DeleteAdRequest{
					AdId:     1,
					AuthorId: 0,
				},
			},
			want:    &proto.DeleteAdResponse{Success: false},
			wantErr: true,
			adExist: errs.AdNotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeApp := mocks.NewIApp(t)
			fakeApp.
				On("DeleteAd", tt.args.ctx, tt.args.request.AdId, tt.args.request.AuthorId).
				Return(tt.adExist).
				Maybe()
			a := &AdService{
				app: fakeApp,
			}
			got, err := a.DeleteAd(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteAd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteAd() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdService_DeleteUser(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *proto.DeleteUserRequest
	}
	tests := []struct {
		name      string
		args      args
		want      *proto.DeleteUserResponse
		userExist error
		wantErr   bool
	}{
		{
			name: "No Error",
			args: args{
				ctx:     context.Background(),
				request: &proto.DeleteUserRequest{Id: 0},
			},
			want:      &proto.DeleteUserResponse{Success: true},
			wantErr:   false,
			userExist: nil,
		},

		{
			name: "not exist",
			args: args{
				ctx:     context.Background(),
				request: &proto.DeleteUserRequest{Id: 2},
			},
			want:      &proto.DeleteUserResponse{Success: false},
			wantErr:   true,
			userExist: errs.UserNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeApp := mocks.NewIApp(t)
			fakeApp.
				On("DeleteUser", tt.args.ctx, tt.args.request.Id).
				Return(tt.userExist).
				Maybe()
			a := &AdService{
				app: fakeApp,
			}
			got, err := a.DeleteUser(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdService_GetUser(t *testing.T) {
	ident := int64(0)
	type args struct {
		ctx     context.Context
		request *proto.GetUserRequest
	}
	tests := []struct {
		name      string
		args      args
		want      *proto.UserResponse
		userExist error
		wantErr   bool
	}{
		{
			name: "No Error",
			args: args{
				ctx:     context.Background(),
				request: &proto.GetUserRequest{Id: &ident},
			},
			want:      &proto.UserResponse{Id: int64(0)},
			wantErr:   false,
			userExist: nil,
		},

		{
			name: "wrong usage",
			args: args{
				ctx:     context.Background(),
				request: &proto.GetUserRequest{Id: nil},
			},
			want:      nil,
			wantErr:   true,
			userExist: errs.WrongProtoBufDataError,
		},

		{
			name: "user not exist",
			args: args{
				ctx:     context.Background(),
				request: &proto.GetUserRequest{Id: &ident},
			},
			want:      nil,
			wantErr:   true,
			userExist: errs.UserNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeApp := mocks.NewIApp(t)
			fakeApp.
				On("FindUser", tt.args.ctx, func() int64 {
					if tt.args.request.Id != nil {
						return *tt.args.request.Id
					}
					return -1
				}()).
				Return(func() *users.User {
					if tt.args.request.Id != nil && *tt.args.request.Id != int64(-1) {
						return &users.User{
							ID: ident,
						}
					}
					return nil
				}(), tt.userExist).
				Maybe()
			a := &AdService{
				app: fakeApp,
			}
			got, err := a.GetUser(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdService_ListAds(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *proto.ListAdRequest
	}
	tests := []struct {
		name     string
		args     args
		want     *proto.ListAdResponse
		adsExist error
		wantErr  bool
		ret      []*ads.Ad
	}{
		{
			name: "No Error",
			args: args{
				ctx:     context.Background(),
				request: &proto.ListAdRequest{Title: "Main"},
			},
			want: &proto.ListAdResponse{List: []*proto.AdResponse{{
				Id:        0,
				Title:     "H",
				Text:      "E",
				AuthorId:  0,
				Published: false,
			},
				{
					Id:        1,
					Title:     "O",
					Text:      "F",
					AuthorId:  0,
					Published: false,
				},
			}},
			wantErr:  false,
			adsExist: nil,
			ret: []*ads.Ad{{
				ID:        0,
				AuthorID:  0,
				Title:     "H",
				Text:      "E",
				Published: false,
			},
				{
					ID:        1,
					AuthorID:  0,
					Title:     "O",
					Text:      "F",
					Published: false,
				},
			},
		},

		{
			name: "not found",
			args: args{
				ctx:     context.Background(),
				request: &proto.ListAdRequest{Title: "Mama"},
			},
			want:     nil,
			ret:      []*ads.Ad{},
			wantErr:  true,
			adsExist: errs.AdNotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeApp := mocks.NewIApp(t)
			fakeApp.
				On("GetAdByName", tt.args.ctx, tt.args.request.Title).
				Return(tt.ret, tt.adsExist).
				Maybe()
			a := &AdService{
				app: fakeApp,
			}
			got, err := a.ListAds(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListAds() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdService_UpdateAd(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *proto.UpdateAdRequest
	}
	tests := []struct {
		name      string
		args      args
		want      *proto.AdResponse
		userExist error
		adExist   error
		wantErr   bool
	}{
		{
			name: "No Error",
			args: args{
				ctx: context.Background(),
				request: &proto.UpdateAdRequest{
					AdId:   0,
					Title:  "hello",
					Text:   "World",
					UserId: 0,
				},
			},
			want: &proto.AdResponse{
				Id:        0,
				Title:     "hello",
				Text:      "World",
				AuthorId:  0,
				Published: false,
			},
			wantErr:   false,
			userExist: nil,
			adExist:   nil,
		},

		{
			name: "Error: no such user",
			args: args{
				ctx: context.Background(),
				request: &proto.UpdateAdRequest{
					AdId:   0,
					Title:  "hello",
					Text:   "World",
					UserId: 1,
				},
			},
			want:      nil,
			wantErr:   true,
			userExist: errs.UserNotFoundError,
			adExist:   nil,
		},

		{
			name: "Error: validation",
			args: args{
				ctx: context.Background(),
				request: &proto.UpdateAdRequest{
					AdId:   0,
					Title:  "",
					Text:   "",
					UserId: 0,
				},
			},
			want:      nil,
			wantErr:   true,
			userExist: nil,
			adExist:   errs.ValidationError,
		},

		{
			name: "Error: permission denied",
			args: args{
				ctx: context.Background(),
				request: &proto.UpdateAdRequest{
					AdId:   0,
					Title:  "hello",
					Text:   "World",
					UserId: 2,
				},
			},
			want:      nil,
			wantErr:   true,
			userExist: nil,
			adExist:   errs.AccessError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeApp := mocks.NewIApp(t)
			fakeApp.
				On("FindUser", tt.args.ctx, tt.args.request.UserId).
				Return(nil, tt.userExist).
				Maybe()
			fakeApp.
				On("UpdateAd", tt.args.ctx, tt.args.request.AdId, tt.args.request.UserId,
					tt.args.request.Title, tt.args.request.Text).
				Return(&ads.Ad{
					ID:        0,
					Title:     tt.args.request.Title,
					Text:      tt.args.request.Text,
					AuthorID:  0,
					CDate:     time.Time{},
					UDate:     time.Time{},
					Published: false,
				}, tt.adExist).
				Maybe()
			a := &AdService{
				app: fakeApp,
			}
			got, err := a.UpdateAd(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateAd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateAd() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdService_UpdateUser(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *proto.UpdateUserRequest
	}
	tests := []struct {
		name      string
		args      args
		want      *proto.UserResponse
		userExist error
		wantErr   bool
	}{
		{
			name: "No Error",
			args: args{
				ctx: context.Background(),
				request: &proto.UpdateUserRequest{
					Id:    0,
					Name:  "N",
					Email: "J",
				},
			},
			want: &proto.UserResponse{
				Id:    0,
				Name:  "N",
				Email: "J",
			},
			wantErr:   false,
			userExist: nil,
		},

		{
			name: "not exist",
			args: args{
				ctx: context.Background(),
				request: &proto.UpdateUserRequest{
					Id:    3,
					Name:  "Ko",
					Email: "GH",
				},
			},
			want:      nil,
			wantErr:   true,
			userExist: errs.UserNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeApp := mocks.NewIApp(t)
			fakeApp.
				On("UpdateUser", tt.args.ctx, tt.args.request.Id,
					tt.args.request.Name, tt.args.request.Email).
				Return(func() *users.User {
					if tt.userExist != nil {
						return nil
					}
					return &users.User{
						ID:    0,
						Name:  tt.args.request.Name,
						Email: tt.args.request.Email,
					}
				}(), tt.userExist,
				).
				Maybe()
			a := &AdService{
				app: fakeApp,
			}
			got, err := a.UpdateUser(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAdService(t *testing.T) {
	a := app.NewApp(repo.NewAd(), repo.NewUser())
	got := NewAdService(a)

	if !reflect.DeepEqual(NewAdService(a), got) {
		t.Errorf("Wrong service creation")
	}
}
