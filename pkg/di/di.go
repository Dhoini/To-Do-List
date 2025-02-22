package di

import (
	"ToDo/internal/models"
	"context"
)

type INoteRepository interface {
	Create(ctx context.Context, note *models.Note) (*models.Note, error)
	GetAll(ctx context.Context, userID string, limit, offset int) ([]models.Note, int64, error)
	Get(ctx context.Context, noteID string) (*models.Note, error)
	Update(ctx context.Context, note *models.Note) (*models.Note, error)
	Delete(ctx context.Context, noteID string) error
}

type INoteService interface {
	CreateNote(ctx context.Context, note *models.Note) (*models.Note, error)
	GetAllNotes(ctx context.Context, userID string, limit, offset int) ([]models.Note, int64, error)
	GetNote(ctx context.Context, noteID string) (*models.Note, error)
	UpdateNote(ctx context.Context, note *models.Note) (*models.Note, error)
	DeleteNote(ctx context.Context, noteID string) error
}

type IAuthService interface {
	Register(ctx context.Context, email, password, name string) (string, error)
	Login(ctx context.Context, email, password string) (*models.User, error)
}

type IUserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	FindById(ctx context.Context, userId string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}
