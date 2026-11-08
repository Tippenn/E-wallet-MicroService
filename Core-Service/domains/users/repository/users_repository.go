package repository

import (
	"context"

	"e-wallet-system-api/domains/users"
	"e-wallet-system-api/domains/users/entities"
	"e-wallet-system-api/infrastructure/database"
)

type userRepository struct {
	db database.Database
}

func NewUserRepository(db database.Database) users.UserRepository {
	return userRepository{db: db}
}

// CreateUser implements [users.UserRepository].
func (u userRepository) CreateUser(ctx context.Context, userId string, email string, name string, password string) error {
	user := &entities.User{
		Id:       userId,
		Email:    email,
		Name:     name,
		Password: password,
		Enabled:  1,
	}

	result := u.db.GetInstance().Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FindUserbyID implements [users.UserRepository].
func (u userRepository) FindUserbyID(Id string) (*entities.User, error) {
	var user *entities.User
	result := u.db.GetInstance().First(&user, "id=?", Id)

	if result.Error != nil {
		println("repo")
		return nil, result.Error
	}

	return user, nil
}

// FindUserbyUsername implements [users.UserRepository].
func (u userRepository) FindUserbyEmail(email string) (*entities.User, error) {
	var user *entities.User
	result := u.db.GetInstance().First(&user, "email=?", email)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// UpdateUser implements [users.UserRepository].
func (u userRepository) UpdateUser(username *entities.User) error {
	panic("unimplemented")
}
