package auth

import (
	"ToDo/internal/user"
	"ToDo/pkg/depInj"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type AuthService struct {
	UserRepository depInj.IUserRepository
}

func NewUserService(userRepository depInj.IUserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New(ErrUserExisted)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	_, err = service.UserRepository.Create(user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

func (service *AuthService) Login(email, password string) (*user.User, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser == nil {
		return nil, errors.New(ErrWrongCredentials)
	}
	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		slog.Info(err.Error())
		return nil, errors.New(ErrWrongCredentials)
	}
	return existedUser, nil
}
