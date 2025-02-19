package main

import (
	"ToDo/configs"
	"ToDo/internal/auth"
	"ToDo/internal/notes"
	"ToDo/internal/user"
	"ToDo/pkg/middleware"

	//"ToDo/internal/user"
	"ToDo/pkg/db"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	Db := db.NewDb(conf)
	router := http.NewServeMux()

	UserRepository := user.NewUserRepository(Db)
	NoteRepository := notes.NewNoteRepository(Db)

	AuthService := auth.NewUserService(UserRepository)

	notes.NewNoteHandler(router, &notes.NoteHandlerDeps{
		NoteRepository: NoteRepository,
		Config:         conf,
	})

	auth.NewAuthHandler(router, &auth.AuthHandlerDeps{
		AuthService: AuthService,
		Config:      conf,
	})

	stackMiddlewares := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	return stackMiddlewares(router)
}
