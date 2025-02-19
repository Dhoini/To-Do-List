package notes

import (
	"ToDo/pkg/db"
	"ToDo/pkg/idgen"
	"ToDo/pkg/valid"
	"errors"
	"gorm.io/gorm"
	"log/slog"
)

type NoteRepository struct {
	DataBase *db.Db
}

func NewNoteRepository(dataBase *db.Db) *NoteRepository {
	return &NoteRepository{
		DataBase: dataBase,
	}
}

func (repo *NoteRepository) Create(note *Note) (*Note, error) {
	note.ID = idgen.GenerateNanoID()
	if note.ID == "" {
		return nil, errors.New("id is empty")
	}
	result := repo.DataBase.Create(note)
	if result.Error != nil {
		slog.Error("Can not create note", result.Error)
		return nil, result.Error
	}

	return note, nil
}

func (repo *NoteRepository) GetAllNotes(limit, offset int) ([]Note, error) {
	var notes []Note
	query := repo.DataBase.
		Table("notes").
		Select("*").
		Session(&gorm.Session{})

	query = query.
		Order("created_at asc").
		Limit(int(limit)).
		Offset(int(offset)).
		Scan(&notes)

	if query.Error != nil {
		return nil, query.Error
	}
	return notes, nil
}

func (repo *NoteRepository) GetById(noteId string) (*Note, error) {
	var note Note
	result := repo.DataBase.Where("id = ?", noteId).First(&note)
	if result.Error != nil {
		return nil, result.Error
	}
	return &note, nil
}

func (repo *NoteRepository) Update(note *Note) (*Note, error) {
	err := valid.ValidateStatus(note.Status)
	if err != nil {
		slog.Error("Invalid note status", note.Status)
		return nil, err
	}

	result := repo.DataBase.Save(note)
	if result.Error != nil {
		slog.Error("Can not update note", result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		slog.Error("Note does not exist", note.ID)
		return nil, gorm.ErrRecordNotFound
	}

	return note, nil
}
func (repo *NoteRepository) Delete(noteId string) error {
	result := repo.DataBase.Where("id = ?", noteId).Delete(&Note{})
	if result.Error != nil {
		slog.Error("Can not delete note", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		slog.Error("Note does not exist", noteId)
		return gorm.ErrRecordNotFound
	}
	return nil
}
