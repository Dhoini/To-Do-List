package notes

import "time"

type CreateNoteRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

type GetAllNotesResponse struct {
	Notes      []Note `json:"notes"`
	TotalCount int64  `json:"total_count"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

type GetNoteResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status" validate:"omitempty,oneof=created in_progress done"`
}
