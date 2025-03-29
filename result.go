package fall

import (
	"encoding/json"
	"net/http"
)

type Result[T any] struct {
	Value T
	Error error
}

func NewResult[T any](value T, err error) Result[T] {
	return Result[T]{value, err}
}

func (r *Result[T]) Unwrap() (T, error) {
	return r.Value, r.Error
}

func (r *Result[T]) UnwrapOr(defaultValue T) T {
	if r.Error != nil {
		return defaultValue
	}
	return r.Value
}

func (r *Result[T]) UnwrapOrElse(f func(error) T) T {
	if r.Error != nil {
		return f(r.Error)
	}
	return r.Value
}

func (r *Result[T]) UnwrapPanic() T {
	if r.Error != nil {
		panic(r.Error)
	}
	return r.Value
}

type Resulter interface {
	Err() error
}

func (r Result[T]) Err() error {
	return r.Error
}

func ValidateResults(results ...Resulter) bool {
	for _, result := range results {
		if result.Err() != nil {
			return false
		}
	}
	return true
}

type ErrorsDTO struct {
	Errors []string `json:"errors"`
}

func RequestValidation(w http.ResponseWriter, r *http.Request, results ...Resulter) bool {
	errorMessages := make([]string, 0)
	for _, result := range results {
		if err := result.Err(); err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}

	if len(errorMessages) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		dto := ErrorsDTO{Errors: errorMessages}
		if err := json.NewEncoder(w).Encode(dto); err != nil {
			http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		}

		return false
	}
	return true
}
