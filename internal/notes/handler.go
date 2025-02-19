package notes

import (
	"ToDo/configs"
	"ToDo/pkg/middleware"
	"ToDo/pkg/req"
	"ToDo/pkg/res"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type NoteHandlerDeps struct {
	NoteRepository *NoteRepository
	Config         *configs.Config
}

type NoteHandler struct {
	NoteRepository *NoteRepository
}

func NewNoteHandler(router *http.ServeMux, deps *NoteHandlerDeps) {
	handler := &NoteHandler{
		NoteRepository: deps.NoteRepository,
	}
	middlewares := middleware.Chain(
		middleware.RateLimiter(1, 3, time.Minute),
		middleware.IsAuthenticated(deps.Config),
	)

	router.Handle("POST /notes", middlewares(handler.CreateNote()))
	router.Handle("GET /notes", middlewares(handler.GetAllNotes()))
	router.Handle("GET /notes/{id}", middlewares(handler.GetNote()))
	router.Handle("PATCH /notes/{id}", middlewares(handler.UpdateNote()))
	router.Handle("DELETE /notes/{id}", middlewares(handler.DeleteNote()))
}

func (handler *NoteHandler) CreateNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := req.HandleBody[Note](&w, r)
		if err != nil {
			return
		}

		createdNote, err := handler.NoteRepository.Create(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.JsonResponse(w, createdNote, http.StatusOK)
	}
}

func (handler *NoteHandler) GetAllNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		links, err := handler.NoteRepository.GetAllNotes(limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		res.JsonResponse(w, links, http.StatusOK)
	}
}

func (handler *NoteHandler) GetNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteId := r.PathValue("id")

		note, err := handler.NoteRepository.GetById(noteId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		res.JsonResponse(w, note, http.StatusOK)

	}

}

//	func (handler *NoteHandler) ChangeStatus() http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			noteId := r.PathValue("id")
//			note, err := handler.NoteRepository.ChangeStatusById(noteId)
//			if err != nil {
//				http.Error(w, err.Error(), http.StatusBadRequest)
//			}
//		}
//	}
func (handler *NoteHandler) UpdateNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[Note](&w, r)
		if err != nil {
			slog.Error("can not handle request")
			return
		}

		noteId := r.PathValue("id")
		if err != nil {
			slog.Error("can not parse note id")
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		note, err := handler.NoteRepository.GetById(noteId)
		if err != nil {
			slog.Error("note not found", "error", err)
			http.Error(w, "note not found", http.StatusNotFound)
			return
		}

		note.Title = body.Title
		note.Content = body.Content
		note.Status = body.Status

		updatedNote, err := handler.NoteRepository.Update(note)
		if err != nil {
			slog.Error("can not update note", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.JsonResponse(w, updatedNote, http.StatusOK)
	}

}

func (handler *NoteHandler) DeleteNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteId := r.PathValue("id")

		_, err := handler.NoteRepository.GetById(noteId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.NoteRepository.Delete(noteId)
		if err != nil {
			slog.Error("can not delete note", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.JsonResponse(w, "deleted", http.StatusOK)

	}
}
