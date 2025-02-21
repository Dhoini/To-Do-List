package main

import (
	"ToDo/configs"
	"ToDo/internal/auth"
	"ToDo/internal/notes"
	"ToDo/internal/user"
	"ToDo/pkg/middleware"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"ToDo/pkg/db"
)

// Константы для улучшения читаемости
const (
	shutdownTimeout = 15 * time.Second
)

func main() {
	// Инициализация конфигурации
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err) // Fatal уже завершает программу
	}

	// Настройка логгера
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	// Подключение к БД
	gormDB, sqlDB, err := db.NewDb(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer closeDB(sqlDB)

	// Миграции
	if err := runMigrations(gormDB); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Инициализация зависимостей
	router := setupRouter(gormDB, cfg)

	// Настройка сервера
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Запуск сервера в горутине
	go startServer(server)

	// Ожидание сигнала завершения
	awaitShutdown(server)
}

// Вспомогательные функции для улучшения читаемости
func closeDB(sqlDB *sql.DB) {
	if err := sqlDB.Close(); err != nil {
		slog.Error("Failed to close database connection", "error", err)
	}
}

func setupRouter(gormDB *gorm.DB, cfg *configs.Config) http.Handler {
	router := http.NewServeMux()

	userRepo := user.NewUserRepository(gormDB)
	noteRepo := notes.NewNoteRepository(gormDB)
	authSvc := auth.NewUserService(userRepo)
	noteSvc := notes.NewNoteService(noteRepo)

	notes.NewNoteHandler(router, &notes.NoteHandler{
		NoteRepository: noteSvc,
		Config:         cfg,
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

func startServer(server *http.Server) {
	slog.Info("Server listening", "port", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Server failed to start", "error", err)
	}
}

func awaitShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}
	slog.Info("Server exited")
}

func runMigrations(db *gorm.DB) error {
	db.Logger = db.Logger.LogMode(logger.Info)
	return db.AutoMigrate(&user.User{}, &notes.Note{})
}
