package repository

import "uas-go/app/model"

type UserRepository interface {
	GetUserByUsername(username string) (*model.User, string, error)
	GetRoleName(roleID string) (string, error)
	GetPermissionsByRole(roleID string) ([]string, error)
	GetUserByID(id string) (*model.User, error)
}

type UserRepo struct{}

var _ UserRepository = (*UserRepo)(nil)

func (r *UserRepo) GetUserByUsername(username string) (*model.User, string, error) {
	return GetUserByUsername(username)
}

func (r *UserRepo) GetUserByID(id string) (*model.User, error) {
	return GetUserByID(id)
}

func (r *UserRepo) GetRoleName(roleID string) (string, error) {
	return GetRoleName(roleID)
}

func (r *UserRepo) GetPermissionsByRole(roleID string) ([]string, error) {
	return GetPermissionsByRole(roleID)
}
