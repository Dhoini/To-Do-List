package idgen

import (
	"github.com/matoous/go-nanoid/v2"
	"log/slog"
)

func GenerateNanoID() string {
	id, err := gonanoid.New()
	if err != nil {
		slog.Error("Failed to generate NanoID", "error", err)
		return ""
	}
	return id
}
