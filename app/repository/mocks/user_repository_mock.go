package mocks

import "uas-go/app/model"

type UserRepositoryMock struct {
	GetUserByUsernameFn func(username string) (*model.User, string, error)
	GetUserByIDFn       func(id string) (*model.User, error)
	GetRoleNameFn       func(roleID string) (string, error)
	GetPermissionsFn    func(roleID string) ([]string, error)
}

func (m *UserRepositoryMock) GetUserByUsername(username string) (*model.User, string, error) {
	return m.GetUserByUsernameFn(username)
}

func (m *UserRepositoryMock) GetUserByID(id string) (*model.User, error) {
	return m.GetUserByIDFn(id)
}

func (m *UserRepositoryMock) GetRoleName(roleID string) (string, error) {
	return m.GetRoleNameFn(roleID)
}

func (m *UserRepositoryMock) GetPermissionsByRole(roleID string) ([]string, error) {
	return m.GetPermissionsFn(roleID)
}
