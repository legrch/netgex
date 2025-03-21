// Code generated by mockery v2.53.3. DO NOT EDIT.

package service

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// Registrar is an autogenerated mock type for the Registrar type
type Registrar struct {
	mock.Mock
}

type Registrar_Expecter struct {
	mock *mock.Mock
}

func (_m *Registrar) EXPECT() *Registrar_Expecter {
	return &Registrar_Expecter{mock: &_m.Mock}
}

// RegisterGRPC provides a mock function with given fields: _a0
func (_m *Registrar) RegisterGRPC(_a0 *grpc.Server) {
	_m.Called(_a0)
}

// Registrar_RegisterGRPC_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterGRPC'
type Registrar_RegisterGRPC_Call struct {
	*mock.Call
}

// RegisterGRPC is a helper method to define mock.On call
//   - _a0 *grpc.Server
func (_e *Registrar_Expecter) RegisterGRPC(_a0 interface{}) *Registrar_RegisterGRPC_Call {
	return &Registrar_RegisterGRPC_Call{Call: _e.mock.On("RegisterGRPC", _a0)}
}

func (_c *Registrar_RegisterGRPC_Call) Run(run func(_a0 *grpc.Server)) *Registrar_RegisterGRPC_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*grpc.Server))
	})
	return _c
}

func (_c *Registrar_RegisterGRPC_Call) Return() *Registrar_RegisterGRPC_Call {
	_c.Call.Return()
	return _c
}

func (_c *Registrar_RegisterGRPC_Call) RunAndReturn(run func(*grpc.Server)) *Registrar_RegisterGRPC_Call {
	_c.Run(run)
	return _c
}

// RegisterHTTP provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Registrar) RegisterHTTP(_a0 context.Context, _a1 *runtime.ServeMux, _a2 string, _a3 []grpc.DialOption) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	if len(ret) == 0 {
		panic("no return value specified for RegisterHTTP")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Registrar_RegisterHTTP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterHTTP'
type Registrar_RegisterHTTP_Call struct {
	*mock.Call
}

// RegisterHTTP is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *runtime.ServeMux
//   - _a2 string
//   - _a3 []grpc.DialOption
func (_e *Registrar_Expecter) RegisterHTTP(_a0 interface{}, _a1 interface{}, _a2 interface{}, _a3 interface{}) *Registrar_RegisterHTTP_Call {
	return &Registrar_RegisterHTTP_Call{Call: _e.mock.On("RegisterHTTP", _a0, _a1, _a2, _a3)}
}

func (_c *Registrar_RegisterHTTP_Call) Run(run func(_a0 context.Context, _a1 *runtime.ServeMux, _a2 string, _a3 []grpc.DialOption)) *Registrar_RegisterHTTP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*runtime.ServeMux), args[2].(string), args[3].([]grpc.DialOption))
	})
	return _c
}

func (_c *Registrar_RegisterHTTP_Call) Return(_a0 error) *Registrar_RegisterHTTP_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Registrar_RegisterHTTP_Call) RunAndReturn(run func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error) *Registrar_RegisterHTTP_Call {
	_c.Call.Return(run)
	return _c
}

// NewRegistrar creates a new instance of Registrar. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRegistrar(t interface {
	mock.TestingT
	Cleanup(func())
}) *Registrar {
	mock := &Registrar{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
