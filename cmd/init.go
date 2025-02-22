package main

import (
	"ToDo/configs"
	"ToDo/internal/auth"
	"ToDo/internal/models"
	"ToDo/internal/notes"
	"ToDo/internal/user"
	"ToDo/pkg/db"
	"ToDo/pkg/middleware"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// App содержит все зависимости приложения
type App struct {
	Router  http.Handler
	SQLDB   *sql.DB
	GormDB  *gorm.DB
	Config  *configs.Config
	Cleanup func() // Функция для закрытия ресурсов
}

// InitializeApp инициализирует приложение и возвращает структуру App
func InitializeApp() (*App, error) {
	// Загружаем конфигурацию
	cfg, err := configs.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Настраиваем логгер
	logLevel := slog.LevelDebug // Можно сделать конфигурируемым через cfg
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))

	// Подключаемся к базе данных
	gormDB, sqlDB, err := db.NewDb(cfg)
	if err != nil {
		return nil, err
	}

	// Выполняем миграции
	if err := runMigrations(gormDB); err != nil {
		sqlDB.Close() // Закрываем соединение в случае ошибки
		return nil, err
	}

	// Инициализируем зависимости и маршрутизатор
	router := setupRouter(gormDB, cfg)

	// Функция для очистки (закрытие базы данных)
	cleanup := func() {
		if err := sqlDB.Close(); err != nil {
			slog.Error("Failed to close database connection", "error", err)
		}
	}

	return &App{
		Router:  router,
		SQLDB:   sqlDB,
		GormDB:  gormDB,
		Config:  cfg,
		Cleanup: cleanup,
	}, nil
}

// runMigrations выполняет миграции базы данных
func runMigrations(db *gorm.DB) error {
	db.Logger = db.Logger.LogMode(logger.Info)
	return db.AutoMigrate(&models.User{}, &models.Note{})
}

// setupRouter инициализирует маршрутизатор с зависимостями
func setupRouter(gormDB *gorm.DB, cfg *configs.Config) http.Handler {
	router := http.NewServeMux()

	userRepo := user.NewUserRepository(gormDB)
	noteRepo := notes.NewNoteRepository(gormDB)
	authSvc := auth.NewUserService(userRepo)
	noteSvc := notes.NewNoteService(noteRepo)

	notes.NewNoteHandler(router, &notes.NoteHandlerDeps{
		NoteService: noteSvc,
		Config:      cfg,
	})
	auth.NewAuthHandler(router, &auth.AuthHandlerDeps{
		AuthService: authSvc,
		Config:      cfg,
	})

	return middleware.Chain(
		middleware.CORS,
		middleware.Logging,
		middleware.RateLimiter(cfg.RateLimit.MaxRequests, cfg.RateLimit.Burst, cfg.RateLimit.TTL),
	)(router)
}
