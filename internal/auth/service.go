package auth

import (
	"ToDo/internal/models"
	"ToDo/internal/user"
	"ToDo/pkg/di"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type AuthService struct {
	UserRepository di.IUserRepository
}

func NewUserService(userRepository di.IUserRepository) *AuthService {
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
	newUser := &models.User{
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

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, error) {
	existingUser, err := s.UserRepository.FindByEmail(ctx, email)
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
