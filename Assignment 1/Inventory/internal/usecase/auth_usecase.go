package usecase

import (
	"errors"
	"inventory/internal/domain"
)

type AuthUseCase interface {
	Register(username, password string) error
	Login(username, password string) (UserResponse, error)
}

type authUseCase struct {
	repo domain.UserRepository
}

func NewAuthUseCase(r domain.UserRepository) AuthUseCase {
	return &authUseCase{repo: r}
}

func (uc *authUseCase) Register(username, password string) error {
	user := domain.User{
		Username: username,
		Password: password,
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	return uc.repo.Create(user)
}

func (uc *authUseCase) Login(username, password string) (UserResponse, error) {
	user, err := uc.repo.GetByUsername(username)
	if err != nil {
		return UserResponse{}, err
	}

	if !user.CheckPassword(password) {
		return UserResponse{}, errors.New("invalid credentials")
	}

	return UserResponse{Username: user.Username}, nil
}

type UserResponse struct {
	Username string `json:"username"`
}
