package notes

import (
	"context"
	"testing"
	"time"

	"ToDo/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNoteRepository — мок для INoteRepository
type MockNoteRepository struct {
	mock.Mock
}

func (m *MockNoteRepository) Create(ctx context.Context, note *models.Note) (*models.Note, error) {
	args := m.Called(ctx, note)
	return args.Get(0).(*models.Note), args.Error(1)
}

func (m *MockNoteRepository) GetAll(ctx context.Context, userID string, limit, offset int) ([]models.Note, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]models.Note), args.Get(1).(int64), args.Error(2)
}

func (m *MockNoteRepository) Get(ctx context.Context, noteID string) (*models.Note, error) {
	args := m.Called(ctx, noteID)
	return args.Get(0).(*models.Note), args.Error(1)
}

func (m *MockNoteRepository) Update(ctx context.Context, note *models.Note) (*models.Note, error) {
	args := m.Called(ctx, note)
	return args.Get(0).(*models.Note), args.Error(1)
}

func (m *MockNoteRepository) Delete(ctx context.Context, noteID string) error {
	args := m.Called(ctx, noteID)
	return args.Error(0)
}

func TestNoteService_CreateNote(t *testing.T) {
	tests := []struct {
		name      string
		note      *models.Note
		mockSetup func(m *MockNoteRepository)
		wantErr   bool
		err       error
		wantNote  *models.Note
	}{
		{
			name: "Successful note creation with default status",
			note: &models.Note{
				Title:   "Test Note",
				Content: "Test content",
				UserID:  "user123",
			},
			mockSetup: func(m *MockNoteRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(note *models.Note) bool {
					return note.Status == "created" && note.Title == "Test Note"
				})).Return(&models.Note{
					ID:        "note123",
					Title:     "Test Note",
					Content:   "Test content",
					Status:    "created",
					UserID:    "user123",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			wantErr:  false,
			wantNote: &models.Note{ID: "note123", Title: "Test Note", Content: "Test content", Status: "created", UserID: "user123"},
		},
		{
			name: "Successful note creation with valid status",
			note: &models.Note{
				Title:   "Test Note",
				Content: "Test content",
				Status:  "in_progress",
				UserID:  "user123",
			},
			mockSetup: func(m *MockNoteRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(note *models.Note) bool {
					return note.Status == "in_progress" && note.Title == "Test Note"
				})).Return(&models.Note{
					ID:        "note123",
					Title:     "Test Note",
					Content:   "Test content",
					Status:    "in_progress",
					UserID:    "user123",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			wantErr:  false,
			wantNote: &models.Note{ID: "note123", Title: "Test Note", Content: "Test content", Status: "in_progress", UserID: "user123"},
		},
		{
			name: "Invalid status returns error",
			note: &models.Note{
				Title:   "Test Note",
				Content: "Test content",
				Status:  "invalid",
				UserID:  "user123",
			},
			mockSetup: func(m *MockNoteRepository) {},
			wantErr:   true,
			err:       ErrInvalidNoteStatus,
		},
		{
			name: "Repository error on create",
			note: &models.Note{
				Title:   "Test Note",
				Content: "Test content",
				UserID:  "user123",
			},
			mockSetup: func(m *MockNoteRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return((*models.Note)(nil), assert.AnError)
			},
			wantErr: true,
			err:     assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок-репозиторий
			mockRepo := new(MockNoteRepository)
			tt.mockSetup(mockRepo)

			// Создаем сервис с мок-репозиторием
			service := NewNoteService(mockRepo)

			// Вызываем метод CreateNote
			ctx := context.Background()
			gotNote, err := service.CreateNote(ctx, tt.note)

			// Проверяем ошибку
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err, "expected error")
				assert.Nil(t, gotNote, "note should be nil on error")
			} else {
				assert.NoError(t, err, "unexpected error")
				assert.NotNil(t, gotNote, "note should not be nil")
				assert.Equal(t, tt.wantNote.ID, gotNote.ID, "note ID mismatch")
				assert.Equal(t, tt.wantNote.Title, gotNote.Title, "note title mismatch")
				assert.Equal(t, tt.wantNote.Content, gotNote.Content, "note content mismatch")
				assert.Equal(t, tt.wantNote.Status, gotNote.Status, "note status mismatch")
				assert.Equal(t, tt.wantNote.UserID, gotNote.UserID, "note userID mismatch")
			}

			// Проверяем, что все ожидаемые вызовы мока были выполнены
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestNoteService_GetAllNotes — тесты для GetAllNotes
func TestNoteService_GetAllNotes(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		limit     int
		offset    int
		mockSetup func(m *MockNoteRepository)
		wantNotes []models.Note
		wantCount int64
		wantErr   bool
		err       error
	}{
		{
			name:   "Successful fetch of notes",
			userID: "user123",
			limit:  10,
			offset: 0,
			mockSetup: func(m *MockNoteRepository) {
				m.On("GetAll", mock.Anything, "user123", 10, 0).Return([]models.Note{
					{ID: "note1", UserID: "user123", Title: "Note 1"},
					{ID: "note2", UserID: "user123", Title: "Note 2"},
				}, int64(2), nil)
			},
			wantNotes: []models.Note{
				{ID: "note1", UserID: "user123", Title: "Note 1"},
				{ID: "note2", UserID: "user123", Title: "Note 2"},
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:   "Repository error on fetch",
			userID: "user123",
			limit:  10,
			offset: 0,
			mockSetup: func(m *MockNoteRepository) {
				m.On("GetAll", mock.Anything, "user123", 10, 0).Return(nil, int64(0), assert.AnError)
			},
			wantNotes: nil,
			wantCount: 0,
			wantErr:   true,
			err:       assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок-репозиторий
			mockRepo := new(MockNoteRepository)
			tt.mockSetup(mockRepo)

			// Создаем сервис с мок-репозиторием
			service := NewNoteService(mockRepo)

			// Вызываем метод GetAllNotes
			ctx := context.Background()
			notes, count, err := service.GetAllNotes(ctx, tt.userID, tt.limit, tt.offset)

			// Проверяем ошибку
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err, "expected error")
				assert.Nil(t, notes, "notes should be nil on error")
				assert.Equal(t, int64(0), count, "count should be 0 on error")
			} else {
				assert.NoError(t, err, "unexpected error")
				assert.NotNil(t, notes, "notes should not be nil")
				assert.Equal(t, tt.wantNotes, notes, "notes mismatch")
				assert.Equal(t, tt.wantCount, count, "count mismatch")
			}

			// Проверяем, что все ожидаемые вызовы мока были выполнены
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestNoteService_GetNote — тесты для GetNote
func TestNoteService_GetNote(t *testing.T) {
	tests := []struct {
		name      string
		noteID    string
		mockSetup func(m *MockNoteRepository)
		wantNote  *models.Note
		wantErr   bool
		err       error
	}{
		{
			name:   "Successful note fetch",
			noteID: "note123",
			mockSetup: func(m *MockNoteRepository) {
				m.On("Get", mock.Anything, "note123").Return(&models.Note{
					ID:     "note123",
					Title:  "Test Note",
					UserID: "user123",
				}, nil)
			},
			wantNote: &models.Note{ID: "note123", Title: "Test Note", UserID: "user123"},
			wantErr:  false,
		},
		{
			name:   "Repository error on fetch",
			noteID: "note123",
			mockSetup: func(m *MockNoteRepository) {
				m.On("Get", mock.Anything, "note123").Return(nil, assert.AnError)
			},
			wantNote: nil,
			wantErr:  true,
			err:      assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок-репозиторий
			mockRepo := new(MockNoteRepository)
			tt.mockSetup(mockRepo)

			// Создаем сервис с мок-репозиторием
			service := NewNoteService(mockRepo)

			// Вызываем метод GetNote
			ctx := context.Background()
			note, err := service.GetNote(ctx, tt.noteID)

			// Проверяем ошибку
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err, "expected error")
				assert.Nil(t, note, "note should be nil on error")
			} else {
				assert.NoError(t, err, "unexpected error")
				assert.NotNil(t, note, "note should not be nil")
				assert.Equal(t, tt.wantNote, note, "note mismatch")
			}

			// Проверяем, что все ожидаемые вызовы мока были выполнены
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestNoteService_UpdateNote — тесты для UpdateNote
func TestNoteService_UpdateNote(t *testing.T) {
	tests := []struct {
		name      string
		note      *models.Note
		mockSetup func(m *MockNoteRepository)
		wantErr   bool
		err       error
		wantNote  *models.Note
	}{
		{
			name: "Successful note update with valid status",
			note: &models.Note{
				ID:     "note123",
				Title:  "Updated Note",
				Status: "done",
				UserID: "user123",
			},
			mockSetup: func(m *MockNoteRepository) {
				m.On("Update", mock.Anything, mock.MatchedBy(func(note *models.Note) bool {
					return note.ID == "note123" && note.Status == "done"
				})).Return(&models.Note{
					ID:        "note123",
					Title:     "Updated Note",
					Status:    "done",
					UserID:    "user123",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			wantErr:  false,
			wantNote: &models.Note{ID: "note123", Title: "Updated Note", Status: "done", UserID: "user123"},
		},
		{
			name: "Invalid status returns error",
			note: &models.Note{
				ID:     "note123",
				Title:  "Updated Note",
				Status: "invalid",
				UserID: "user123",
			},
			mockSetup: func(m *MockNoteRepository) {},
			wantErr:   true,
			err:       ErrInvalidNoteStatus,
		},
		{
			name: "Repository error on update",
			note: &models.Note{
				ID:     "note123",
				Title:  "Updated Note",
				UserID: "user123",
			},
			mockSetup: func(m *MockNoteRepository) {
				m.On("Update", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			wantErr: true,
			err:     assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок-репозиторий
			mockRepo := new(MockNoteRepository)
			tt.mockSetup(mockRepo)

			// Создаем сервис с мок-репозиторием
			service := NewNoteService(mockRepo)

			// Вызываем метод UpdateNote
			ctx := context.Background()
			gotNote, err := service.UpdateNote(ctx, tt.note)

			// Проверяем ошибку
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err, "expected error")
				assert.Nil(t, gotNote, "note should be nil on error")
			} else {
				assert.NoError(t, err, "unexpected error")
				assert.NotNil(t, gotNote, "note should not be nil")
				assert.Equal(t, tt.wantNote.ID, gotNote.ID, "note ID mismatch")
				assert.Equal(t, tt.wantNote.Title, gotNote.Title, "note title mismatch")
				assert.Equal(t, tt.wantNote.Status, gotNote.Status, "note status mismatch")
				assert.Equal(t, tt.wantNote.UserID, gotNote.UserID, "note userID mismatch")
			}

			// Проверяем, что все ожидаемые вызовы мока были выполнены
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestNoteService_DeleteNote — тесты для DeleteNote
func TestNoteService_DeleteNote(t *testing.T) {
	tests := []struct {
		name      string
		noteID    string
		mockSetup func(m *MockNoteRepository)
		wantErr   bool
		err       error
	}{
		{
			name:   "Successful note deletion",
			noteID: "note123",
			mockSetup: func(m *MockNoteRepository) {
				m.On("Delete", mock.Anything, "note123").Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Repository error on deletion",
			noteID: "note123",
			mockSetup: func(m *MockNoteRepository) {
				m.On("Delete", mock.Anything, "note123").Return(assert.AnError)
			},
			wantErr: true,
			err:     assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок-репозиторий
			mockRepo := new(MockNoteRepository)
			tt.mockSetup(mockRepo)

			// Создаем сервис с мок-репозиторием
			service := NewNoteService(mockRepo)

			// Вызываем метод DeleteNote
			ctx := context.Background()
			err := service.DeleteNote(ctx, tt.noteID)

			// Проверяем ошибку
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err, "expected error")
			} else {
				assert.NoError(t, err, "unexpected error")
			}

			// Проверяем, что все ожидаемые вызовы мока были выполнены
			mockRepo.AssertExpectations(t)
		})
	}
}
