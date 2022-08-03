// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/mongo (interfaces: Collection,Query)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	mongo "github.com/juju/juju/mongo"
	mgo "github.com/juju/mgo/v3"
)

// MockCollection is a mock of Collection interface.
type MockCollection struct {
	ctrl     *gomock.Controller
	recorder *MockCollectionMockRecorder
}

// MockCollectionMockRecorder is the mock recorder for MockCollection.
type MockCollectionMockRecorder struct {
	mock *MockCollection
}

// NewMockCollection creates a new mock instance.
func NewMockCollection(ctrl *gomock.Controller) *MockCollection {
	mock := &MockCollection{ctrl: ctrl}
	mock.recorder = &MockCollectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCollection) EXPECT() *MockCollectionMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockCollection) Count() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockCollectionMockRecorder) Count() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockCollection)(nil).Count))
}

// Find mocks base method.
func (m *MockCollection) Find(arg0 interface{}) mongo.Query {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0)
	ret0, _ := ret[0].(mongo.Query)
	return ret0
}

// Find indicates an expected call of Find.
func (mr *MockCollectionMockRecorder) Find(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockCollection)(nil).Find), arg0)
}

// FindId mocks base method.
func (m *MockCollection) FindId(arg0 interface{}) mongo.Query {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindId", arg0)
	ret0, _ := ret[0].(mongo.Query)
	return ret0
}

// FindId indicates an expected call of FindId.
func (mr *MockCollectionMockRecorder) FindId(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindId", reflect.TypeOf((*MockCollection)(nil).FindId), arg0)
}

// Name mocks base method.
func (m *MockCollection) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockCollectionMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockCollection)(nil).Name))
}

// Pipe mocks base method.
func (m *MockCollection) Pipe(arg0 interface{}) *mgo.Pipe {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pipe", arg0)
	ret0, _ := ret[0].(*mgo.Pipe)
	return ret0
}

// Pipe indicates an expected call of Pipe.
func (mr *MockCollectionMockRecorder) Pipe(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pipe", reflect.TypeOf((*MockCollection)(nil).Pipe), arg0)
}

// Writeable mocks base method.
func (m *MockCollection) Writeable() mongo.WriteCollection {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Writeable")
	ret0, _ := ret[0].(mongo.WriteCollection)
	return ret0
}

// Writeable indicates an expected call of Writeable.
func (mr *MockCollectionMockRecorder) Writeable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Writeable", reflect.TypeOf((*MockCollection)(nil).Writeable))
}

// MockQuery is a mock of Query interface.
type MockQuery struct {
	ctrl     *gomock.Controller
	recorder *MockQueryMockRecorder
}

// MockQueryMockRecorder is the mock recorder for MockQuery.
type MockQueryMockRecorder struct {
	mock *MockQuery
}

// NewMockQuery creates a new mock instance.
func NewMockQuery(ctrl *gomock.Controller) *MockQuery {
	mock := &MockQuery{ctrl: ctrl}
	mock.recorder = &MockQueryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuery) EXPECT() *MockQueryMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockQuery) All(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// All indicates an expected call of All.
func (mr *MockQueryMockRecorder) All(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockQuery)(nil).All), arg0)
}

// Apply mocks base method.
func (m *MockQuery) Apply(arg0 mgo.Change, arg1 interface{}) (*mgo.ChangeInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Apply", arg0, arg1)
	ret0, _ := ret[0].(*mgo.ChangeInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Apply indicates an expected call of Apply.
func (mr *MockQueryMockRecorder) Apply(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Apply", reflect.TypeOf((*MockQuery)(nil).Apply), arg0, arg1)
}

// Batch mocks base method.
func (m *MockQuery) Batch(arg0 int) mongo.Query {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Batch", arg0)
	ret0, _ := ret[0].(mongo.Query)
	return ret0
}

// Batch indicates an expected call of Batch.
func (mr *MockQueryMockRecorder) Batch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Batch", reflect.TypeOf((*MockQuery)(nil).Batch), arg0)
}

// Comment mocks base method.
func (m *MockQuery) Comment(arg0 string) mongo.Query {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Comment", arg0)
	ret0, _ := ret[0].(mongo.Query)
	return ret0
}

// Comment indicates an expected call of Comment.
func (mr *MockQueryMockRecorder) Comment(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Comment", reflect.TypeOf((*MockQuery)(nil).Comment), arg0)
}

// Count mocks base method.
func (m *MockQuery) Count() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockQueryMockRecorder) Count() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockQuery)(nil).Count))
}

// Distinct mocks base method.
func (m *MockQuery) Distinct(arg0 string, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Distinct", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Distinct indicates an expected call of Distinct.
func (mr *MockQueryMockRecorder) Distinct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Distinct", reflect.TypeOf((*MockQuery)(nil).Distinct), arg0, arg1)
}

// Explain mocks base method.
func (m *MockQuery) Explain(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Explain", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Explain indicates an expected call of Explain.
func (mr *MockQueryMockRecorder) Explain(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Explain", reflect.TypeOf((*MockQuery)(nil).Explain), arg0)
}

// For mocks base method.
func (m *MockQuery) For(arg0 interface{}, arg1 func() error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "For", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// For indicates an expected call of For.
func (mr *MockQueryMockRecorder) For(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "For", reflect.TypeOf((*MockQuery)(nil).For), arg0, arg1)
}

// Hint mocks base method.
func (m *MockQuery) Hint(arg0 ...string) mongo.Query {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Hint", varargs...)
	ret0, _ := ret[0].(mongo.Query)
	return ret0
}

// Hint indicates an expected call of Hint.
func (mr *MockQueryMockRecorder) Hint(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hint", reflect.TypeOf((*MockQuery)(nil).Hint), arg0...)
}

// Iter mocks base method.
func (m *MockQuery) Iter() mongo.Iterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Iter")
	ret0, _ := ret[0].(mongo.Iterator)
	return ret0
}

// Iter indicates an expected call of Iter.
func (mr *MockQueryMockRecorder) Iter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Iter", reflect.TypeOf((*MockQuery)(nil).Iter))
}

// Limit mocks base method.
func (m *MockQuery) Limit(arg0 int) mongo.Query {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Limit", arg0)
	ret0, _ := ret[0].(mongo.Query)
	return ret0
}

// Limit indicates an expected call of Limit.
func (mr *MockQueryMockRecorder) Limit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Limit", reflect.TypeOf((*MockQuery)(nil).Limit), arg0)
}

// LogReplay mocks base method.
func (m *MockQuery) LogReplay() mongo.Query {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogReplay")
	ret0, _ := ret[0].(mongo.Query)
	return ret0
}

// LogReplay indicates an expected call of LogReplay.
func (mr *MockQueryMockRecorder) LogReplay() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogReplay", reflect.TypeOf((*MockQuery)(nil).LogReplay))
}

// MapReduce mocks base method.
func (m *MockQuery) MapReduce(arg0 *mgo.MapReduce, arg1 interface{}) (*mgo.MapReduceInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MapReduce", arg0, arg1)
	ret0, _ := ret[0].(*mgo.MapReduceInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MapReduce indicates an expected call of MapReduce.
func (mr *MockQueryMockRecorder) MapReduce(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MapReduce", reflect.TypeOf((*MockQuery)(nil).MapReduce), arg0, arg1)
}

// One mocks base method.
func (m *MockQuery) One(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "One", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// One indicates an expected call of One.
func (mr *MockQueryMockRecorder) One(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "One", reflect.TypeOf((*MockQuery)(nil).One), arg0)
}

// Prefetch mocks base method.
func (m *MockQuery) Prefetch(arg0 float64) mongo.Query {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prefetch", arg0)
	ret0, _ := ret[0].(mongo.Query)
	return ret0
}

// Prefetch indicates an expected call of Prefetch.
func (mr *MockQueryMockRecorder) Prefetch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prefetch", reflect.TypeOf((*MockQuery)(nil).Prefetch), arg0)
}

// Select mocks base method.
func (m *MockQuery) Select(arg0 interface{}) mongo.Query {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Select", arg0)
	ret0, _ := ret[0].(mongo.Query)
	return ret0
}

// Select indicates an expected call of Select.
func (mr *MockQueryMockRecorder) Select(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Select", reflect.TypeOf((*MockQuery)(nil).Select), arg0)
}

// SetMaxScan mocks base method.
func (m *MockQuery) SetMaxScan(arg0 int) mongo.Query {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetMaxScan", arg0)
	ret0, _ := ret[0].(mongo.Query)
	return ret0
}

// SetMaxScan indicates an expected call of SetMaxScan.
func (mr *MockQueryMockRecorder) SetMaxScan(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxScan", reflect.TypeOf((*MockQuery)(nil).SetMaxScan), arg0)
}

// SetMaxTime mocks base method.
func (m *MockQuery) SetMaxTime(arg0 time.Duration) mongo.Query {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetMaxTime", arg0)
	ret0, _ := ret[0].(mongo.Query)
	return ret0
}

// SetMaxTime indicates an expected call of SetMaxTime.
func (mr *MockQueryMockRecorder) SetMaxTime(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxTime", reflect.TypeOf((*MockQuery)(nil).SetMaxTime), arg0)
}

// Skip mocks base method.
func (m *MockQuery) Skip(arg0 int) mongo.Query {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Skip", arg0)
	ret0, _ := ret[0].(mongo.Query)
	return ret0
}

// Skip indicates an expected call of Skip.
func (mr *MockQueryMockRecorder) Skip(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Skip", reflect.TypeOf((*MockQuery)(nil).Skip), arg0)
}

// Snapshot mocks base method.
func (m *MockQuery) Snapshot() mongo.Query {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Snapshot")
	ret0, _ := ret[0].(mongo.Query)
	return ret0
}

// Snapshot indicates an expected call of Snapshot.
func (mr *MockQueryMockRecorder) Snapshot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Snapshot", reflect.TypeOf((*MockQuery)(nil).Snapshot))
}

// Sort mocks base method.
func (m *MockQuery) Sort(arg0 ...string) mongo.Query {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Sort", varargs...)
	ret0, _ := ret[0].(mongo.Query)
	return ret0
}

// Sort indicates an expected call of Sort.
func (mr *MockQueryMockRecorder) Sort(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sort", reflect.TypeOf((*MockQuery)(nil).Sort), arg0...)
}

// Tail mocks base method.
func (m *MockQuery) Tail(arg0 time.Duration) *mgo.Iter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tail", arg0)
	ret0, _ := ret[0].(*mgo.Iter)
	return ret0
}

// Tail indicates an expected call of Tail.
func (mr *MockQueryMockRecorder) Tail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tail", reflect.TypeOf((*MockQuery)(nil).Tail), arg0)
}
