package notes

import (
	"ToDo/pkg/idgen"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type NoteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(dataBase *gorm.DB) *NoteRepository {
	return &NoteRepository{
		db: dataBase,
	}
}

func (r *NoteRepository) CreateNote(ctx context.Context, note *Note) (*Note, error) {
	note.ID = idgen.GenerateNanoID()
	if note.ID == "" {
		return nil, fmt.Errorf("generate id: %w", ErrCreateNote)
	}

	result := r.db.WithContext(ctx).Create(note)
	if result.Error != nil {
		return nil, fmt.Errorf("create note with ID %s: %w", note.ID, result.Error)
	}
	return note, nil
}

func (r *NoteRepository) GetAllNotes(ctx context.Context, userId string, limit, offset int) ([]Note, int64, error) {
	var notes []Note
	var totalCount int64

	countQuery := r.db.WithContext(ctx).Model(&Note{}).Where("user_id = ?", userId).Count(&totalCount)
	if countQuery.Error != nil {
		return nil, 0, fmt.Errorf("get total count for user %s: %w", userId, countQuery.Error)
	}

	query := r.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Order("created_at asc").
		Limit(limit).
		Offset(offset).
		Find(&notes) // Исправлено ¬es на &notes

	if query.Error != nil {
		return nil, 0, fmt.Errorf("get all notes for user %s: %w", userId, query.Error)
	}
	return notes, totalCount, nil
}

func (r *NoteRepository) GetNote(ctx context.Context, noteId string) (*Note, error) {
	var note Note
	result := r.db.WithContext(ctx).Where("id = ?", noteId).First(&note) // Исправлено ¬e на &note
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get note by id %s: %w", noteId, ErrNoteNotFound)
		}
		return nil, fmt.Errorf("get note by id %s: %w", noteId, result.Error)
	}
	return &note, nil
}

func (r *NoteRepository) UpdateNote(ctx context.Context, note *Note) (*Note, error) {
	result := r.db.WithContext(ctx).Save(note)
	if result.Error != nil {
		return nil, fmt.Errorf("update note with ID %s: %w", note.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("update note with ID %s: %w", note.ID, ErrNoteNotFound)
	}
	return note, nil
}

func (r *NoteRepository) DeleteNote(ctx context.Context, noteId string) error {
	result := r.db.WithContext(ctx).Where("id = ?", noteId).Delete(&Note{})
	if result.Error != nil {
		return fmt.Errorf("delete note with ID %s: %w", noteId, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("delete note with ID %s: %w", noteId, ErrNoteNotFound)
	}
	return nil
}
