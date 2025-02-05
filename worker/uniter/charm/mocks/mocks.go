// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/worker/uniter/charm (interfaces: BundleReader,BundleInfo,Bundle)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	set "github.com/juju/collections/set"
	charm "github.com/juju/juju/worker/uniter/charm"
	gomock "go.uber.org/mock/gomock"
)

// MockBundleReader is a mock of BundleReader interface.
type MockBundleReader struct {
	ctrl     *gomock.Controller
	recorder *MockBundleReaderMockRecorder
}

// MockBundleReaderMockRecorder is the mock recorder for MockBundleReader.
type MockBundleReaderMockRecorder struct {
	mock *MockBundleReader
}

// NewMockBundleReader creates a new mock instance.
func NewMockBundleReader(ctrl *gomock.Controller) *MockBundleReader {
	mock := &MockBundleReader{ctrl: ctrl}
	mock.recorder = &MockBundleReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBundleReader) EXPECT() *MockBundleReaderMockRecorder {
	return m.recorder
}

// Read mocks base method.
func (m *MockBundleReader) Read(arg0 charm.BundleInfo, arg1 <-chan struct{}) (charm.Bundle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0, arg1)
	ret0, _ := ret[0].(charm.Bundle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockBundleReaderMockRecorder) Read(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockBundleReader)(nil).Read), arg0, arg1)
}

// MockBundleInfo is a mock of BundleInfo interface.
type MockBundleInfo struct {
	ctrl     *gomock.Controller
	recorder *MockBundleInfoMockRecorder
}

// MockBundleInfoMockRecorder is the mock recorder for MockBundleInfo.
type MockBundleInfoMockRecorder struct {
	mock *MockBundleInfo
}

// NewMockBundleInfo creates a new mock instance.
func NewMockBundleInfo(ctrl *gomock.Controller) *MockBundleInfo {
	mock := &MockBundleInfo{ctrl: ctrl}
	mock.recorder = &MockBundleInfoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBundleInfo) EXPECT() *MockBundleInfoMockRecorder {
	return m.recorder
}

// ArchiveSha256 mocks base method.
func (m *MockBundleInfo) ArchiveSha256() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArchiveSha256")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ArchiveSha256 indicates an expected call of ArchiveSha256.
func (mr *MockBundleInfoMockRecorder) ArchiveSha256() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArchiveSha256", reflect.TypeOf((*MockBundleInfo)(nil).ArchiveSha256))
}

// String mocks base method.
func (m *MockBundleInfo) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockBundleInfoMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockBundleInfo)(nil).String))
}

// MockBundle is a mock of Bundle interface.
type MockBundle struct {
	ctrl     *gomock.Controller
	recorder *MockBundleMockRecorder
}

// MockBundleMockRecorder is the mock recorder for MockBundle.
type MockBundleMockRecorder struct {
	mock *MockBundle
}

// NewMockBundle creates a new mock instance.
func NewMockBundle(ctrl *gomock.Controller) *MockBundle {
	mock := &MockBundle{ctrl: ctrl}
	mock.recorder = &MockBundleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBundle) EXPECT() *MockBundleMockRecorder {
	return m.recorder
}

// ArchiveMembers mocks base method.
func (m *MockBundle) ArchiveMembers() (set.Strings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArchiveMembers")
	ret0, _ := ret[0].(set.Strings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ArchiveMembers indicates an expected call of ArchiveMembers.
func (mr *MockBundleMockRecorder) ArchiveMembers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArchiveMembers", reflect.TypeOf((*MockBundle)(nil).ArchiveMembers))
}

// ExpandTo mocks base method.
func (m *MockBundle) ExpandTo(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExpandTo", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExpandTo indicates an expected call of ExpandTo.
func (mr *MockBundleMockRecorder) ExpandTo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExpandTo", reflect.TypeOf((*MockBundle)(nil).ExpandTo), arg0)
}
