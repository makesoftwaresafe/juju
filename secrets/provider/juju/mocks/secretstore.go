// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/state (interfaces: SecretsStore)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	secrets "github.com/juju/juju/core/secrets"
	state "github.com/juju/juju/state"
)

// MockSecretsStore is a mock of SecretsStore interface.
type MockSecretsStore struct {
	ctrl     *gomock.Controller
	recorder *MockSecretsStoreMockRecorder
}

// MockSecretsStoreMockRecorder is the mock recorder for MockSecretsStore.
type MockSecretsStoreMockRecorder struct {
	mock *MockSecretsStore
}

// NewMockSecretsStore creates a new mock instance.
func NewMockSecretsStore(ctrl *gomock.Controller) *MockSecretsStore {
	mock := &MockSecretsStore{ctrl: ctrl}
	mock.recorder = &MockSecretsStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretsStore) EXPECT() *MockSecretsStoreMockRecorder {
	return m.recorder
}

// CreateSecret mocks base method.
func (m *MockSecretsStore) CreateSecret(arg0 *secrets.URI, arg1 state.CreateSecretParams) (*secrets.SecretMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSecret", arg0, arg1)
	ret0, _ := ret[0].(*secrets.SecretMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSecret indicates an expected call of CreateSecret.
func (mr *MockSecretsStoreMockRecorder) CreateSecret(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSecret", reflect.TypeOf((*MockSecretsStore)(nil).CreateSecret), arg0, arg1)
}

// DeleteSecret mocks base method.
func (m *MockSecretsStore) DeleteSecret(arg0 *secrets.URI) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecret", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSecret indicates an expected call of DeleteSecret.
func (mr *MockSecretsStoreMockRecorder) DeleteSecret(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecret", reflect.TypeOf((*MockSecretsStore)(nil).DeleteSecret), arg0)
}

// GetSecret mocks base method.
func (m *MockSecretsStore) GetSecret(arg0 *secrets.URI) (*secrets.SecretMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", arg0)
	ret0, _ := ret[0].(*secrets.SecretMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret.
func (mr *MockSecretsStoreMockRecorder) GetSecret(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockSecretsStore)(nil).GetSecret), arg0)
}

// GetSecretValue mocks base method.
func (m *MockSecretsStore) GetSecretValue(arg0 *secrets.URI, arg1 int) (secrets.SecretValue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecretValue", arg0, arg1)
	ret0, _ := ret[0].(secrets.SecretValue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecretValue indicates an expected call of GetSecretValue.
func (mr *MockSecretsStoreMockRecorder) GetSecretValue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretValue", reflect.TypeOf((*MockSecretsStore)(nil).GetSecretValue), arg0, arg1)
}

// ListSecrets mocks base method.
func (m *MockSecretsStore) ListSecrets(arg0 state.SecretsFilter) ([]*secrets.SecretMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecrets", arg0)
	ret0, _ := ret[0].([]*secrets.SecretMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecrets indicates an expected call of ListSecrets.
func (mr *MockSecretsStoreMockRecorder) ListSecrets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecrets", reflect.TypeOf((*MockSecretsStore)(nil).ListSecrets), arg0)
}

// UpdateSecret mocks base method.
func (m *MockSecretsStore) UpdateSecret(arg0 *secrets.URI, arg1 state.UpdateSecretParams) (*secrets.SecretMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSecret", arg0, arg1)
	ret0, _ := ret[0].(*secrets.SecretMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSecret indicates an expected call of UpdateSecret.
func (mr *MockSecretsStoreMockRecorder) UpdateSecret(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSecret", reflect.TypeOf((*MockSecretsStore)(nil).UpdateSecret), arg0, arg1)
}
