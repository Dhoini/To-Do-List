package auth

import (
	"ToDo/internal/user"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type AuthService struct {
	UserRepository IUserRepository
}

type IUserRepository interface {
	Create(ctx context.Context, user *user.User) (*user.User, error)
	FindById(ctx context.Context, userId string) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}

func NewUserService(userRepository IUserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (string, error) {
	//Проверяем существует ли пользователь (используем FindByEmail)
	_, err := s.UserRepository.FindByEmail(ctx, email)
	if err == nil { // Если ошибки нет, значит, пользователь найден.
		return "", user.ErrUserAlreadyExists
	}
	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return "", fmt.Errorf("check user existance: %w", err)
	}
	//Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	//Создаем пользователя
	newUser := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	createdUser, err := s.UserRepository.Create(ctx, newUser)
	if err != nil {
		return "", err // Ошибка уже обернута в репозитории
	}
	return createdUser.ID, nil
}

func (s *AuthService) Login(ctx context.Context, id, password string) (*user.User, error) {
	existingUser, err := s.UserRepository.FindById(ctx, id)
	if err != nil {
		return nil, err // Ошибка уже обернута в репозитории
	}
	// Сравниваем хешированный пароль
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))
	if err != nil {
		slog.Info("Invalid password", "error", err)
		return nil, user.ErrUserNotFound // Возвращаем ErrUserNotFound для безопасности
	}
	return existingUser, nil
}
