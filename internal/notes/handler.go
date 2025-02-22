package notes

import (
	"ToDo/configs"
	"ToDo/pkg/di"
	"ToDo/pkg/middleware"
	"net/http"
)

type NoteHandlerDeps struct {
	Config      *configs.Config
	NoteService di.INoteService
}
type NoteHandler struct {
	Config      *configs.Config
	NoteService di.INoteService
}

func NewNoteHandler(router *http.ServeMux, deps *NoteHandlerDeps) {
	handler := &NoteHandler{
		Config:      deps.Config,
		NoteService: deps.NoteService,
	}
	middlewares := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
		middleware.RateLimiter(deps.Config.RateLimit.MaxRequests, deps.Config.RateLimit.Burst, deps.Config.RateLimit.TTL),
		middleware.IsAuthenticated(deps.Config),
	)

	router.Handle("POST /notes", middlewares(handler.CreateNote()))
	router.Handle("GET /notes", middlewares(handler.GetAllNotes()))
	router.Handle("GET /notes/{id}", middlewares(handler.GetNote()))
	router.Handle("PATCH /notes/{id}", middlewares(handler.UpdateNote()))
	router.Handle("DELETE /notes/{id}", middlewares(handler.DeleteNote()))
}
