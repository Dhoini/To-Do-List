package notes

import (
	"ToDo/internal/models"
	"ToDo/pkg/di"
	"context"
	"log/slog"
)

type NoteService struct {
	noteRepository di.INoteRepository // Используем интерфейс вместо конкретного типа
}

func NewNoteService(noteRepo di.INoteRepository) *NoteService { // Принимаем интерфейс
	return &NoteService{noteRepository: noteRepo}
}

func (s *NoteService) CreateNote(ctx context.Context, note *models.Note) (*models.Note, error) {
	validStatuses := map[string]bool{"created": true, "in_progress": true, "done": true}
	if note.Status == "" {
		note.Status = "created" // Значение по умолчанию
	} else if !validStatuses[note.Status] {
		return nil, ErrInvalidNoteStatus
	}

	slog.Info("Creating note", "title", note.Title, "user_id", note.UserID)
	return s.noteRepository.Create(ctx, note)
}

func (s *NoteService) GetAllNotes(ctx context.Context, userID string, limit, offset int) ([]models.Note, int64, error) {
	slog.Info("Fetching all notes", "user_id", userID, "limit", limit, "offset", offset)
	return s.noteRepository.GetAll(ctx, userID, limit, offset)
}

func (s *NoteService) GetNote(ctx context.Context, noteID string) (*models.Note, error) {
	slog.Info("Fetching note", "note_id", noteID)
	return s.noteRepository.Get(ctx, noteID)
}

func (s *NoteService) UpdateNote(ctx context.Context, note *models.Note) (*models.Note, error) {
	validStatuses := map[string]bool{"created": true, "in_progress": true, "done": true}
	if note.Status != "" && !validStatuses[note.Status] {
		return nil, ErrInvalidNoteStatus
	}
	return s.noteRepository.Update(ctx, note)
}

func (s *NoteService) DeleteNote(ctx context.Context, noteID string) error {
	slog.Info("Deleting note", "note_id", noteID)
	return s.noteRepository.Delete(ctx, noteID)
}
