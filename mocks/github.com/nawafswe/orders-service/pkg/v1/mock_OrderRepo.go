// Code generated by mockery v2.40.1. DO NOT EDIT.

package v1

import (
	context "context"

	models "github.com/nawafswe/orders-service/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MockOrderRepo is an autogenerated mock type for the OrderRepo type
type MockOrderRepo struct {
	mock.Mock
}

type MockOrderRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *MockOrderRepo) EXPECT() *MockOrderRepo_Expecter {
	return &MockOrderRepo_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, order
func (_m *MockOrderRepo) Create(ctx context.Context, order models.Order) (models.Order, error) {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Order) (models.Order, error)); ok {
		return rf(ctx, order)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.Order) models.Order); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Get(0).(models.Order)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.Order) error); ok {
		r1 = rf(ctx, order)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrderRepo_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockOrderRepo_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - order models.Order
func (_e *MockOrderRepo_Expecter) Create(ctx interface{}, order interface{}) *MockOrderRepo_Create_Call {
	return &MockOrderRepo_Create_Call{Call: _e.mock.On("Create", ctx, order)}
}

func (_c *MockOrderRepo_Create_Call) Run(run func(ctx context.Context, order models.Order)) *MockOrderRepo_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.Order))
	})
	return _c
}

func (_c *MockOrderRepo_Create_Call) Return(_a0 models.Order, _a1 error) *MockOrderRepo_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrderRepo_Create_Call) RunAndReturn(run func(context.Context, models.Order) (models.Order, error)) *MockOrderRepo_Create_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateOrderStatus provides a mock function with given fields: ctx, id, status
func (_m *MockOrderRepo) UpdateOrderStatus(ctx context.Context, id int64, status string) (models.Order, error) {
	ret := _m.Called(ctx, id, status)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrderStatus")
	}

	var r0 models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) (models.Order, error)); ok {
		return rf(ctx, id, status)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) models.Order); ok {
		r0 = rf(ctx, id, status)
	} else {
		r0 = ret.Get(0).(models.Order)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, string) error); ok {
		r1 = rf(ctx, id, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrderRepo_UpdateOrderStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateOrderStatus'
type MockOrderRepo_UpdateOrderStatus_Call struct {
	*mock.Call
}

// UpdateOrderStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
//   - status string
func (_e *MockOrderRepo_Expecter) UpdateOrderStatus(ctx interface{}, id interface{}, status interface{}) *MockOrderRepo_UpdateOrderStatus_Call {
	return &MockOrderRepo_UpdateOrderStatus_Call{Call: _e.mock.On("UpdateOrderStatus", ctx, id, status)}
}

func (_c *MockOrderRepo_UpdateOrderStatus_Call) Run(run func(ctx context.Context, id int64, status string)) *MockOrderRepo_UpdateOrderStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(string))
	})
	return _c
}

func (_c *MockOrderRepo_UpdateOrderStatus_Call) Return(_a0 models.Order, _a1 error) *MockOrderRepo_UpdateOrderStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrderRepo_UpdateOrderStatus_Call) RunAndReturn(run func(context.Context, int64, string) (models.Order, error)) *MockOrderRepo_UpdateOrderStatus_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockOrderRepo creates a new instance of MockOrderRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOrderRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOrderRepo {
	mock := &MockOrderRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
