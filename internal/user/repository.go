package user

import (
	"ToDo/pkg/idgen"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(dataBase *gorm.DB) *UserRepository {
	return &UserRepository{db: dataBase}
}

func (r *UserRepository) Create(ctx context.Context, user *User) (*User, error) {
	user.ID = idgen.GenerateNanoID()
	if user.ID == "" {
		return nil, fmt.Errorf("generate id: %w", errors.New("failed to generate id")) // Конкретная ошибка
	}

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) { // Проверяем на дубликат
			return nil, fmt.Errorf("create user: %w", ErrUserAlreadyExists)
		}
		return nil, fmt.Errorf("create user: %w", result.Error)
	}
	return user, nil
}

func (r *UserRepository) FindById(ctx context.Context, userId string) (*User, error) {
	var user User
	result := r.db.WithContext(ctx).Where("id = ?", userId).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) { // Проверяем на ErrRecordNotFound
			return nil, fmt.Errorf("find by id: %w", ErrUserNotFound)
		}
		return nil, fmt.Errorf("find by id: %w", result.Error)
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) { // Проверяем на ErrRecordNotFound
			return nil, fmt.Errorf("find by id: %w", ErrUserNotFound)
		}
		return nil, fmt.Errorf("find by id: %w", result.Error)
	}
	return &user, nil
}
