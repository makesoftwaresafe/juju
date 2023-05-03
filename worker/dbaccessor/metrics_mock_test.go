// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/prometheus/client_golang/prometheus (interfaces: Registerer)

// Package dbaccessor is a generated GoMock package.
package dbaccessor

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	prometheus "github.com/prometheus/client_golang/prometheus"
)

// MockRegisterer is a mock of Registerer interface.
type MockRegisterer struct {
	ctrl     *gomock.Controller
	recorder *MockRegistererMockRecorder
}

// MockRegistererMockRecorder is the mock recorder for MockRegisterer.
type MockRegistererMockRecorder struct {
	mock *MockRegisterer
}

// NewMockRegisterer creates a new mock instance.
func NewMockRegisterer(ctrl *gomock.Controller) *MockRegisterer {
	mock := &MockRegisterer{ctrl: ctrl}
	mock.recorder = &MockRegistererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegisterer) EXPECT() *MockRegistererMockRecorder {
	return m.recorder
}

// MustRegister mocks base method.
func (m *MockRegisterer) MustRegister(arg0 ...prometheus.Collector) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "MustRegister", varargs...)
}

// MustRegister indicates an expected call of MustRegister.
func (mr *MockRegistererMockRecorder) MustRegister(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MustRegister", reflect.TypeOf((*MockRegisterer)(nil).MustRegister), arg0...)
}

// Register mocks base method.
func (m *MockRegisterer) Register(arg0 prometheus.Collector) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockRegistererMockRecorder) Register(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockRegisterer)(nil).Register), arg0)
}

// Unregister mocks base method.
func (m *MockRegisterer) Unregister(arg0 prometheus.Collector) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unregister", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Unregister indicates an expected call of Unregister.
func (mr *MockRegistererMockRecorder) Unregister(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unregister", reflect.TypeOf((*MockRegisterer)(nil).Unregister), arg0)
}
