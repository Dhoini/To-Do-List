package req

import (
	"ToDo/pkg/res"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

func Decode[T any](body io.ReadCloser) (T, error) {
	var payload T
	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func IsValid[T any](payload T) error {
	validate := validator.New()
	err := validate.Struct(payload)
	return err
}

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		res.JsonResponse(*w, err, http.StatusUnprocessableEntity)
		return nil, err
	}

	err = IsValid(body)
	if err != nil {
		res.JsonResponse(*w, err, http.StatusUnprocessableEntity)
		return nil, err
	}
	return &body, nil
}
