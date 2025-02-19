package valid

import (
	"errors"
	"strings"
)

const (
	statusOpen       = "open"
	statusInProgress = "in progress"
	statusReview     = "review"
	statusCompleted  = "completed"
	statusCanceled   = "canceled"
	statusPostponed  = "postponed"
	statusArchived   = "archived"
)

// Список допустимых статусов
var ValidStatuses = []string{
	statusOpen,
	statusInProgress,
	statusReview,
	statusCompleted,
	statusCanceled,
	statusPostponed,
	statusArchived,
}

// ValidateStatus проверяет, является ли переданный статус допустимым
func ValidateStatus(status string) error {
	for _, validStatus := range ValidStatuses {
		if strings.ToLower(status) == validStatus {
			return nil
		}
	}
	return errors.New("invalid status value")
}
