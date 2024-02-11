// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	grpc "ads-server/proto"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IAdService is an autogenerated mock type for the IAdService type
type IAdService struct {
	mock.Mock
}

// ChangeAdStatus provides a mock function with given fields: ctx, request
func (_m *IAdService) ChangeAdStatus(ctx context.Context, request *grpc.ChangeAdStatusRequest) (*grpc.AdResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 *grpc.AdResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.ChangeAdStatusRequest) (*grpc.AdResponse, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.ChangeAdStatusRequest) *grpc.AdResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.AdResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *grpc.ChangeAdStatusRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAd provides a mock function with given fields: ctx, request
func (_m *IAdService) CreateAd(ctx context.Context, request *grpc.CreateAdRequest) (*grpc.AdResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 *grpc.AdResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.CreateAdRequest) (*grpc.AdResponse, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.CreateAdRequest) *grpc.AdResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.AdResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *grpc.CreateAdRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: ctx, request
func (_m *IAdService) CreateUser(ctx context.Context, request *grpc.CreateUserRequest) (*grpc.UserResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 *grpc.UserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.CreateUserRequest) (*grpc.UserResponse, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.CreateUserRequest) *grpc.UserResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.UserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *grpc.CreateUserRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAd provides a mock function with given fields: ctx, request
func (_m *IAdService) DeleteAd(ctx context.Context, request *grpc.DeleteAdRequest) (*grpc.DeleteAdResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 *grpc.DeleteAdResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.DeleteAdRequest) (*grpc.DeleteAdResponse, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.DeleteAdRequest) *grpc.DeleteAdResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.DeleteAdResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *grpc.DeleteAdRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: ctx, request
func (_m *IAdService) GetUser(ctx context.Context, request *grpc.GetUserRequest) (*grpc.UserResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 *grpc.UserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.GetUserRequest) (*grpc.UserResponse, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.GetUserRequest) *grpc.UserResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.UserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *grpc.GetUserRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAds provides a mock function with given fields: ctx, request
func (_m *IAdService) ListAds(ctx context.Context, request *grpc.ListAdRequest) (*grpc.ListAdResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 *grpc.ListAdResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.ListAdRequest) (*grpc.ListAdResponse, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.ListAdRequest) *grpc.ListAdResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.ListAdResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *grpc.ListAdRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAd provides a mock function with given fields: ctx, request
func (_m *IAdService) UpdateAd(ctx context.Context, request *grpc.UpdateAdRequest) (*grpc.AdResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 *grpc.AdResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.UpdateAdRequest) (*grpc.AdResponse, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.UpdateAdRequest) *grpc.AdResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.AdResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *grpc.UpdateAdRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: ctx, request
func (_m *IAdService) UpdateUser(ctx context.Context, request *grpc.UpdateUserRequest) (*grpc.UserResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 *grpc.UserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.UpdateUserRequest) (*grpc.UserResponse, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.UpdateUserRequest) *grpc.UserResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.UserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *grpc.UpdateUserRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIAdService interface {
	mock.TestingT
	Cleanup(func())
}

// NewIAdService creates a new instance of IAdService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIAdService(t mockConstructorTestingTNewIAdService) *IAdService {
	mock := &IAdService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}