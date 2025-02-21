package notes

import (
	"context"
	"log/slog"
)

type NoteService struct {
	noteRepository INoteRepository
}

type INoteRepository interface {
	CreateNote(ctx context.Context, note *Note) (*Note, error)
	GetAllNotes(ctx context.Context, userID string, limit, offset int) ([]Note, int64, error)
	GetNote(ctx context.Context, noteID string) (*Note, error)
	UpdateNote(ctx context.Context, note *Note) (*Note, error)
	DeleteNote(ctx context.Context, noteID string) error
}

func NewNoteService(noteRepo *NoteRepository) *NoteService {
	return &NoteService{noteRepository: noteRepo}
}

func (s *NoteService) Create(ctx context.Context, note *Note) (*Note, error) {
	validStatuses := map[string]bool{"created": true, "in_progress": true, "done": true}
	if note.Status == "" {
		note.Status = "created" // Значение по умолчанию
	} else if !validStatuses[note.Status] {
		return nil, ErrInvalidNoteStatus
	}

	slog.Info("Creating note", "title", note.Title, "user_id", note.UserID)
	return s.noteRepository.CreateNote(ctx, note)
}

func (s *NoteService) GetAll(ctx context.Context, userID string, limit, offset int) ([]Note, int64, error) {
	slog.Info("Fetching all notes", "user_id", userID, "limit", limit, "offset", offset)
	return s.noteRepository.GetAllNotes(ctx, userID, limit, offset)
}

func (s *NoteService) Get(ctx context.Context, noteID string) (*Note, error) {
	slog.Info("Fetching note", "note_id", noteID)
	return s.noteRepository.GetNote(ctx, noteID)
}

func (s *NoteService) Update(ctx context.Context, note *Note) (*Note, error) {
	validStatuses := map[string]bool{"created": true, "in_progress": true, "done": true}
	if note.Status != "" && !validStatuses[note.Status] {
		return nil, ErrInvalidNoteStatus
	}
	return s.noteRepository.UpdateNote(ctx, note)
}

func (s *NoteService) Delete(ctx context.Context, noteID string) error {
	slog.Info("Deleting note", "note_id", noteID)
	return s.noteRepository.DeleteNote(ctx, noteID)
}
