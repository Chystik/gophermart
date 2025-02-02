// Code generated by mockery v2.23.4. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/Chystik/gophermart/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// OrderInteractor is an autogenerated mock type for the OrderInteractor type
type OrderInteractor struct {
	mock.Mock
}

type OrderInteractor_Expecter struct {
	mock *mock.Mock
}

func (_m *OrderInteractor) EXPECT() *OrderInteractor_Expecter {
	return &OrderInteractor_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *OrderInteractor) Create(_a0 context.Context, _a1 models.Order) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Order) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OrderInteractor_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type OrderInteractor_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.Order
func (_e *OrderInteractor_Expecter) Create(_a0 interface{}, _a1 interface{}) *OrderInteractor_Create_Call {
	return &OrderInteractor_Create_Call{Call: _e.mock.On("Create", _a0, _a1)}
}

func (_c *OrderInteractor_Create_Call) Run(run func(_a0 context.Context, _a1 models.Order)) *OrderInteractor_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.Order))
	})
	return _c
}

func (_c *OrderInteractor_Create_Call) Return(_a0 error) *OrderInteractor_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OrderInteractor_Create_Call) RunAndReturn(run func(context.Context, models.Order) error) *OrderInteractor_Create_Call {
	_c.Call.Return(run)
	return _c
}

// GetList provides a mock function with given fields: _a0, _a1
func (_m *OrderInteractor) GetList(_a0 context.Context, _a1 models.User) ([]models.Order, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User) ([]models.Order, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.User) []models.Order); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.User) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderInteractor_GetList_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetList'
type OrderInteractor_GetList_Call struct {
	*mock.Call
}

// GetList is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 models.User
func (_e *OrderInteractor_Expecter) GetList(_a0 interface{}, _a1 interface{}) *OrderInteractor_GetList_Call {
	return &OrderInteractor_GetList_Call{Call: _e.mock.On("GetList", _a0, _a1)}
}

func (_c *OrderInteractor_GetList_Call) Run(run func(_a0 context.Context, _a1 models.User)) *OrderInteractor_GetList_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.User))
	})
	return _c
}

func (_c *OrderInteractor_GetList_Call) Return(_a0 []models.Order, _a1 error) *OrderInteractor_GetList_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrderInteractor_GetList_Call) RunAndReturn(run func(context.Context, models.User) ([]models.Order, error)) *OrderInteractor_GetList_Call {
	_c.Call.Return(run)
	return _c
}

// NewOrderInteractor creates a new instance of OrderInteractor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrderInteractor(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrderInteractor {
	mock := &OrderInteractor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
