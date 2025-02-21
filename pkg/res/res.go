package res

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"` // Сообщение об ошибке (для пользователя)
	//  ErrorCode   int    `json:"error_code,omitempty"`   //  код ошибки (для программной обработки)
	//  Details     any    `json:"details,omitempty"`      //  детали ошибки (например, ошибки валидации)
	//  DebugInfo   string `json:"debug_info,omitempty"` // Отладочная информация (ТОЛЬКО в development среде!)
}

func JsonResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
