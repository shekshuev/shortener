package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) SetURL(key, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockStore) GetURL(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockStore) CheckDBConnection() error {
	args := m.Called()
	return args.Error(0)
}
