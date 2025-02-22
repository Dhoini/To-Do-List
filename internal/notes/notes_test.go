package notes

import (
	"ToDo/internal/models"
	"context"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"testing"
)

type MockNoteService struct {
	mock.Mock
}

func (m *MockNoteService) CreateNote(ctx context.Context, note *models.Note) (*models.Note, error) {
	args := m.Called(ctx, note)
	return args.Get(0).(*models.Note), args.Error(1)
}
func (m *MockNoteService) GetAllNotes(ctx context.Context, noteId string, limit, offset int) ([]models.Note, int64, error) {
	args := m.Called(ctx, noteId, limit, offset)
	return args.Get(0).([]models.Note), args.Get(1).(int64), args.Error(2)
}
func (m *MockNoteService) GetNote(ctx context.Context, noteId string) (*models.Note, error) {
	args := m.Called(ctx, noteId)
	return args.Get(0).(*models.Note), args.Error(1)
}
func (m *MockNoteService) UpdateNote(ctx context.Context, note *models.Note) (*models.Note, error) {
	args := m.Called(ctx, note)
	return args.Get(0).(*models.Note), args.Error(1)
}
func (m *MockNoteService) DeleteNote(ctx context.Context, noteId string) error {
	args := m.Called(ctx, noteId)
	return args.Error(0)
}

func TestNoteHandler_CreateNote(t *testing.T) {
	tests := []struct {
		name               string
		body               CreateNoteRequest
		mockCreate         func(m *MockNoteService)
		expectedStatusCode int
		checkResponse      func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "success create note",
			body: CreateNoteRequest{
				Title: "title",
				Content: "content",
				Status:  "created",
			},
			mockCreate: func(m *MockNoteService) {
				m.On("CreateNote", mock.Anything, mock.Anything,)
			}
		},
	}

}
