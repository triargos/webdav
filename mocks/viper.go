package mocks

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
)

type EnvironmentServiceMock struct {
	mock.Mock
}

func (m *EnvironmentServiceMock) Get(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *EnvironmentServiceMock) GetBool(key string) bool {
	args := m.Called(key)
	return args.Bool(0)
}

type ViperMock struct {
	mock.Mock
}

func (v *ViperMock) SetConfigType(in string) {
	v.Called(in)
}

func (v *ViperMock) SetConfigFile(in string) {
	v.Called(in)
}

func (v *ViperMock) ReadInConfig() error {
	args := v.Called()
	return args.Error(0)
}

func (v *ViperMock) WriteConfig() error {
	args := v.Called()
	return args.Error(0)
}

func (v *ViperMock) Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	args := v.Called(rawVal)
	return args.Error(0)
}

func (v *ViperMock) SetDefault(key string, value interface{}) {
	v.Called(key, value)
}

func (v *ViperMock) Set(key string, value interface{}) {
	v.Called(key, value)
}

func NewViperMock() *ViperMock {
	return &ViperMock{}
}

func NewEnvironmentServiceMock() *EnvironmentServiceMock {
	return &EnvironmentServiceMock{}
}
