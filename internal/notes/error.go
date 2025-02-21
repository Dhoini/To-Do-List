package notes

import "errors"

var (
	ErrNoteNotFound      = errors.New("note not found")
	ErrCreateNote        = errors.New("failed to create note") // и другие
	ErrInvalidNoteStatus = errors.New("invalid note status")
)
