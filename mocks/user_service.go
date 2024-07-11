package mocks

import "github.com/triargos/webdav/pkg/config"

// MockUserService is a mock implementation of the Service interface for testing
type MockUserService struct {
	AddUserFn               func(username string, user config.User) error
	GetUserFn               func(username string) config.User
	GetUsersFn              func() map[string]config.User
	HasUserFn               func(username string) bool
	RemoveUserFn            func(username string) error
	InitializeDirectoriesFn func() error
	HashPasswordsFn         func() error

	AddUserCalls               int
	GetUserCalls               int
	GetUsersCalls              int
	HasUserCalls               int
	RemoveUserCalls            int
	InitializeDirectoriesCalls int
	HashPasswordsCalls         int
}

func NewMockUserService(users map[string]config.User) *MockUserService {

	return &MockUserService{
		AddUserFn:  func(username string, user config.User) error { return nil },
		GetUserFn:  func(username string) config.User { return users[username] },
		GetUsersFn: func() map[string]config.User { return users },
		HasUserFn: func(username string) bool {
			_, ok := users[username]
			return ok
		},
		RemoveUserFn:            func(username string) error { return nil },
		InitializeDirectoriesFn: func() error { return nil },
		HashPasswordsFn:         func() error { return nil },
	}
}

func (m *MockUserService) AddUser(username string, user config.User) error {
	m.AddUserCalls++
	return m.AddUserFn(username, user)
}

func (m *MockUserService) GetUser(username string) config.User {
	m.GetUserCalls++
	return m.GetUserFn(username)
}

func (m *MockUserService) GetUsers() map[string]config.User {
	m.GetUsersCalls++
	return m.GetUsersFn()
}

func (m *MockUserService) HasUser(username string) bool {
	m.HasUserCalls++
	return m.HasUserFn(username)
}

func (m *MockUserService) RemoveUser(username string) error {
	m.RemoveUserCalls++
	return m.RemoveUserFn(username)
}

func (m *MockUserService) InitializeDirectories() error {
	m.InitializeDirectoriesCalls++
	return m.InitializeDirectoriesFn()
}

func (m *MockUserService) HashPasswords() error {
	m.HashPasswordsCalls++
	return m.HashPasswordsFn()
}

// Reset resets all function implementations and call counters
func (m *MockUserService) Reset() {
	m.AddUserFn = func(username string, user config.User) error { return nil }
	m.GetUserFn = func(username string) config.User { return config.User{} }
	m.GetUsersFn = func() map[string]config.User { return make(map[string]config.User) }
	m.HasUserFn = func(username string) bool { return true }
	m.RemoveUserFn = func(username string) error { return nil }
	m.InitializeDirectoriesFn = func() error { return nil }
	m.HashPasswordsFn = func() error { return nil }

	m.AddUserCalls = 0
	m.GetUserCalls = 0
	m.GetUsersCalls = 0
	m.HasUserCalls = 0
	m.RemoveUserCalls = 0
	m.InitializeDirectoriesCalls = 0
	m.HashPasswordsCalls = 0
}
