package notes

import (
	"ToDo/internal/models"
	"ToDo/pkg/middleware"
	"ToDo/pkg/req"
	"ToDo/pkg/res"
	"errors"
	"net/http"
	"strconv"
)

// Вспомогательные функции
func getUserId(r *http.Request) string {
	userId, _ := r.Context().Value(middleware.ContextUserIDKey).(string)
	return userId
}

func parsePagination(r *http.Request) (int, int) {
	const (
		defaultLimit = 10
		maxLimit     = 100
	)

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = defaultLimit
	} else if limit > maxLimit {
		limit = maxLimit
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	return limit, offset
}

func (h *NoteHandler) CreateNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[CreateNoteRequest](&w, r)
		if err != nil {
			return
		}

		userID := getUserId(r)
		if userID == "" {
			res.JsonResponse(w, res.ErrorResponse{Error: "unauthorized"}, http.StatusUnauthorized)
			return
		}

		note := &models.Note{
			Title:   body.Title,
			Content: body.Content,
			Status:  body.Status,
			UserID:  userID,
		}

		createdNote, err := h.NoteService.CreateNote(r.Context(), note)
		if err != nil {
			if errors.Is(err, ErrInvalidNoteStatus) {
				res.JsonResponse(w, res.ErrorResponse{Error: "invalid note status"}, http.StatusBadRequest)
			} else {
				res.JsonResponse(w, res.ErrorResponse{Error: "failed to create note"}, http.StatusInternalServerError)
			}
			return
		}

		res.JsonResponse(w, createdNote, http.StatusCreated)
	}
}

func (h *NoteHandler) GetAllNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := getUserId(r)
		if userId == "" {
			res.JsonResponse(w, res.ErrorResponse{Error: "unauthorized"}, http.StatusUnauthorized)
			return
		}

		limit, offset := parsePagination(r)
		notes, totalCount, err := h.NoteService.GetAllNotes(r.Context(), userId, limit, offset)
		if err != nil {
			res.JsonResponse(w, res.ErrorResponse{Error: "failed to get notes"}, http.StatusInternalServerError)
			return
		}

		res.JsonResponse(w, GetAllNotesResponse{
			Notes:      notes,
			TotalCount: totalCount,
			Limit:      limit,
			Offset:     offset,
		}, http.StatusOK)
	}
}

func (h *NoteHandler) GetNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteId := r.PathValue("id")
		if noteId == "" {
			res.JsonResponse(w, res.ErrorResponse{Error: "note id is required"}, http.StatusBadRequest)
			return
		}

		userId := getUserId(r)
		if userId == "" {
			res.JsonResponse(w, res.ErrorResponse{Error: "unauthorized"}, http.StatusUnauthorized)
			return
		}

		note, err := h.NoteService.GetNote(r.Context(), noteId)
		if err != nil {
			if errors.Is(err, ErrNoteNotFound) {
				res.JsonResponse(w, res.ErrorResponse{Error: "note not found"}, http.StatusNotFound)
			} else {
				res.JsonResponse(w, res.ErrorResponse{Error: "failed to get note by id"}, http.StatusInternalServerError)
			}
			return
		}
		if note.UserID != userId {
			res.JsonResponse(w, res.ErrorResponse{Error: "unauthorized"}, http.StatusUnauthorized)
			return
		}
		response := GetNoteResponse{
			ID:      note.ID,
			Title:   note.Title,
			Content: note.Content,
			Status:  note.Status,
		}

		res.JsonResponse(w, response, http.StatusOK)

	}

}

func (h *NoteHandler) UpdateNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteId := r.PathValue("id")
		if noteId == "" {
			res.JsonResponse(w, res.ErrorResponse{Error: "note id is required"}, http.StatusBadRequest)
			return
		}
		body, err := req.HandleBody[UpdateNoteRequest](&w, r)
		if err != nil {
			return
		}
		userId := getUserId(r)
		if userId == "" {
			res.JsonResponse(w, res.ErrorResponse{Error: "unauthorized"}, http.StatusUnauthorized)
			return
		}
		existingNote, err := h.NoteService.GetNote(r.Context(), noteId)
		if err != nil {
			if errors.Is(err, ErrNoteNotFound) {
				res.JsonResponse(w, res.ErrorResponse{Error: "note not found"}, http.StatusNotFound)
			} else {
				res.JsonResponse(w, res.ErrorResponse{Error: "failed to get note by id"}, http.StatusInternalServerError)
			}
			return
		}
		if existingNote.UserID != userId {
			res.JsonResponse(w, res.ErrorResponse{Error: "unauthorized"}, http.StatusUnauthorized)
			return
		}
		// Обновляем только непустые поля
		if body.Title != "" {
			existingNote.Title = body.Title
		}
		if body.Content != "" {
			existingNote.Content = body.Content
		}
		if body.Status != "" {
			existingNote.Status = body.Status
		}

		updatedNote, err := h.NoteService.UpdateNote(r.Context(), existingNote)
		if err != nil {
			if errors.Is(err, ErrInvalidNoteStatus) {
				res.JsonResponse(w, res.ErrorResponse{Error: "invalid note status"}, http.StatusBadRequest)
			} else {
				res.JsonResponse(w, res.ErrorResponse{Error: "failed to update note"}, http.StatusInternalServerError)
			}
			return
		}
		response := GetNoteResponse{
			ID:        updatedNote.ID,
			Title:     updatedNote.Title,
			Content:   updatedNote.Content,
			Status:    updatedNote.Status,
			UserID:    updatedNote.UserID,
			CreatedAt: updatedNote.CreatedAt,
			UpdatedAt: updatedNote.UpdatedAt,
		}
		res.JsonResponse(w, response, http.StatusOK)
	}
}

func (h *NoteHandler) DeleteNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteId := r.PathValue("id")
		if noteId == "" {
			res.JsonResponse(w, res.ErrorResponse{Error: "note id is required"}, http.StatusBadRequest)
			return
		}

		_, err := h.NoteService.GetNote(r.Context(), noteId)
		if err != nil {
			if errors.Is(err, ErrNoteNotFound) {
				res.JsonResponse(w, res.ErrorResponse{Error: "note not found"}, http.StatusNotFound)
			} else {
				res.JsonResponse(w, res.ErrorResponse{Error: "failed to delete note"}, http.StatusInternalServerError)
			}
			return
		}

		userId := getUserId(r)
		if userId == "" {
			res.JsonResponse(w, res.ErrorResponse{Error: "unauthorized"}, http.StatusUnauthorized)
			return
		}

		err = h.NoteService.DeleteNote(r.Context(), noteId)
		if err != nil {
			res.JsonResponse(w, res.ErrorResponse{Error: "failed to delete note"}, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}
