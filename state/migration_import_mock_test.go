// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/state (interfaces: TransactionRunner,StateDocumentFactory,DocModelNamespace)

// Package state is a generated GoMock package.
package state

import (
	reflect "reflect"

	description "github.com/juju/description/v3"
	txn "github.com/juju/mgo/v2/txn"
	gomock "go.uber.org/mock/gomock"
)

// MockTransactionRunner is a mock of TransactionRunner interface.
type MockTransactionRunner struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionRunnerMockRecorder
}

// MockTransactionRunnerMockRecorder is the mock recorder for MockTransactionRunner.
type MockTransactionRunnerMockRecorder struct {
	mock *MockTransactionRunner
}

// NewMockTransactionRunner creates a new mock instance.
func NewMockTransactionRunner(ctrl *gomock.Controller) *MockTransactionRunner {
	mock := &MockTransactionRunner{ctrl: ctrl}
	mock.recorder = &MockTransactionRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionRunner) EXPECT() *MockTransactionRunnerMockRecorder {
	return m.recorder
}

// RunTransaction mocks base method.
func (m *MockTransactionRunner) RunTransaction(arg0 []txn.Op) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunTransaction", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunTransaction indicates an expected call of RunTransaction.
func (mr *MockTransactionRunnerMockRecorder) RunTransaction(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunTransaction", reflect.TypeOf((*MockTransactionRunner)(nil).RunTransaction), arg0)
}

// MockStateDocumentFactory is a mock of StateDocumentFactory interface.
type MockStateDocumentFactory struct {
	ctrl     *gomock.Controller
	recorder *MockStateDocumentFactoryMockRecorder
}

// MockStateDocumentFactoryMockRecorder is the mock recorder for MockStateDocumentFactory.
type MockStateDocumentFactoryMockRecorder struct {
	mock *MockStateDocumentFactory
}

// NewMockStateDocumentFactory creates a new mock instance.
func NewMockStateDocumentFactory(ctrl *gomock.Controller) *MockStateDocumentFactory {
	mock := &MockStateDocumentFactory{ctrl: ctrl}
	mock.recorder = &MockStateDocumentFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStateDocumentFactory) EXPECT() *MockStateDocumentFactoryMockRecorder {
	return m.recorder
}

// MakeRemoteApplicationDoc mocks base method.
func (m *MockStateDocumentFactory) MakeRemoteApplicationDoc(arg0 description.RemoteApplication) *remoteApplicationDoc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeRemoteApplicationDoc", arg0)
	ret0, _ := ret[0].(*remoteApplicationDoc)
	return ret0
}

// MakeRemoteApplicationDoc indicates an expected call of MakeRemoteApplicationDoc.
func (mr *MockStateDocumentFactoryMockRecorder) MakeRemoteApplicationDoc(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeRemoteApplicationDoc", reflect.TypeOf((*MockStateDocumentFactory)(nil).MakeRemoteApplicationDoc), arg0)
}

// MakeStatusDoc mocks base method.
func (m *MockStateDocumentFactory) MakeStatusDoc(arg0 description.Status) statusDoc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeStatusDoc", arg0)
	ret0, _ := ret[0].(statusDoc)
	return ret0
}

// MakeStatusDoc indicates an expected call of MakeStatusDoc.
func (mr *MockStateDocumentFactoryMockRecorder) MakeStatusDoc(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeStatusDoc", reflect.TypeOf((*MockStateDocumentFactory)(nil).MakeStatusDoc), arg0)
}

// MakeStatusOp mocks base method.
func (m *MockStateDocumentFactory) MakeStatusOp(arg0 string, arg1 statusDoc) txn.Op {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeStatusOp", arg0, arg1)
	ret0, _ := ret[0].(txn.Op)
	return ret0
}

// MakeStatusOp indicates an expected call of MakeStatusOp.
func (mr *MockStateDocumentFactoryMockRecorder) MakeStatusOp(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeStatusOp", reflect.TypeOf((*MockStateDocumentFactory)(nil).MakeStatusOp), arg0, arg1)
}

// NewRemoteApplication mocks base method.
func (m *MockStateDocumentFactory) NewRemoteApplication(arg0 *remoteApplicationDoc) *RemoteApplication {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRemoteApplication", arg0)
	ret0, _ := ret[0].(*RemoteApplication)
	return ret0
}

// NewRemoteApplication indicates an expected call of NewRemoteApplication.
func (mr *MockStateDocumentFactoryMockRecorder) NewRemoteApplication(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRemoteApplication", reflect.TypeOf((*MockStateDocumentFactory)(nil).NewRemoteApplication), arg0)
}

// MockDocModelNamespace is a mock of DocModelNamespace interface.
type MockDocModelNamespace struct {
	ctrl     *gomock.Controller
	recorder *MockDocModelNamespaceMockRecorder
}

// MockDocModelNamespaceMockRecorder is the mock recorder for MockDocModelNamespace.
type MockDocModelNamespaceMockRecorder struct {
	mock *MockDocModelNamespace
}

// NewMockDocModelNamespace creates a new mock instance.
func NewMockDocModelNamespace(ctrl *gomock.Controller) *MockDocModelNamespace {
	mock := &MockDocModelNamespace{ctrl: ctrl}
	mock.recorder = &MockDocModelNamespaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDocModelNamespace) EXPECT() *MockDocModelNamespaceMockRecorder {
	return m.recorder
}

// DocID mocks base method.
func (m *MockDocModelNamespace) DocID(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DocID", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// DocID indicates an expected call of DocID.
func (mr *MockDocModelNamespaceMockRecorder) DocID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DocID", reflect.TypeOf((*MockDocModelNamespace)(nil).DocID), arg0)
}
