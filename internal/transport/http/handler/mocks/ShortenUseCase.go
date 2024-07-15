// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	usecase "github.com/zenoleg/shortener/internal/usecase"
)

// ShortenUseCase is an autogenerated mock type for the ShortenUseCase type
type ShortenUseCase struct {
	mock.Mock
}

// Do provides a mock function with given fields: ctx, query
func (_m *ShortenUseCase) Do(ctx context.Context, query usecase.ShortenQuery) (usecase.DestinationURL, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for Do")
	}

	var r0 usecase.DestinationURL
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usecase.ShortenQuery) (usecase.DestinationURL, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usecase.ShortenQuery) usecase.DestinationURL); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Get(0).(usecase.DestinationURL)
	}

	if rf, ok := ret.Get(1).(func(context.Context, usecase.ShortenQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewShortenUseCase creates a new instance of ShortenUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewShortenUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *ShortenUseCase {
	mock := &ShortenUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
