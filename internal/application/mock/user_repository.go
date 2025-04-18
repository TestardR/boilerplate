// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=../../../application/mock/user_repository.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	model "boilerplate/internal/domain/user/model"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockPersister is a mock of Persister interface.
type MockPersister struct {
	ctrl     *gomock.Controller
	recorder *MockPersisterMockRecorder
	isgomock struct{}
}

// MockPersisterMockRecorder is the mock recorder for MockPersister.
type MockPersisterMockRecorder struct {
	mock *MockPersister
}

// NewMockPersister creates a new mock instance.
func NewMockPersister(ctrl *gomock.Controller) *MockPersister {
	mock := &MockPersister{ctrl: ctrl}
	mock.recorder = &MockPersisterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPersister) EXPECT() *MockPersisterMockRecorder {
	return m.recorder
}

// Persist mocks base method.
func (m *MockPersister) Persist(ctx context.Context, userModel model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Persist", ctx, userModel)
	ret0, _ := ret[0].(error)
	return ret0
}

// Persist indicates an expected call of Persist.
func (mr *MockPersisterMockRecorder) Persist(ctx, userModel any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Persist", reflect.TypeOf((*MockPersister)(nil).Persist), ctx, userModel)
}

// MockLoader is a mock of Loader interface.
type MockLoader struct {
	ctrl     *gomock.Controller
	recorder *MockLoaderMockRecorder
	isgomock struct{}
}

// MockLoaderMockRecorder is the mock recorder for MockLoader.
type MockLoaderMockRecorder struct {
	mock *MockLoader
}

// NewMockLoader creates a new mock instance.
func NewMockLoader(ctrl *gomock.Controller) *MockLoader {
	mock := &MockLoader{ctrl: ctrl}
	mock.recorder = &MockLoaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoader) EXPECT() *MockLoaderMockRecorder {
	return m.recorder
}

// Load mocks base method.
func (m *MockLoader) Load(ctx context.Context, userID model.ID) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", ctx, userID)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load indicates an expected call of Load.
func (mr *MockLoaderMockRecorder) Load(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockLoader)(nil).Load), ctx, userID)
}

// MockUpdater is a mock of Updater interface.
type MockUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockUpdaterMockRecorder
	isgomock struct{}
}

// MockUpdaterMockRecorder is the mock recorder for MockUpdater.
type MockUpdaterMockRecorder struct {
	mock *MockUpdater
}

// NewMockUpdater creates a new mock instance.
func NewMockUpdater(ctrl *gomock.Controller) *MockUpdater {
	mock := &MockUpdater{ctrl: ctrl}
	mock.recorder = &MockUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUpdater) EXPECT() *MockUpdaterMockRecorder {
	return m.recorder
}

// Update mocks base method.
func (m *MockUpdater) Update(ctx context.Context, userModel model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, userModel)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUpdaterMockRecorder) Update(ctx, userModel any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUpdater)(nil).Update), ctx, userModel)
}
