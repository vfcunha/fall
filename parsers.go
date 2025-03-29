package fall

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/mail"
	"reflect"
	"strconv"

	"github.com/google/uuid"
)

func PathValueUUID(r *http.Request, key string) Result[uuid.UUID] {
	return NewResult(uuid.Parse(r.PathValue(key)))
}

func PathValueInt64(r *http.Request, key string) Result[int64] {
	return NewResult(strconv.ParseInt(r.PathValue(key), 10, 64))
}

func PathValueInt(r *http.Request, key string) Result[int] {
	value, err := strconv.ParseInt(r.PathValue(key), 10, 32)
	return NewResult(int(value), err)
}

func PathValue(r *http.Request, key string) Result[string] {
	value := r.PathValue(key)
	if len(value) == 0 {
		return NewResult(value, fmt.Errorf("invalid path value: %s", key))
	}
	return NewResult(value, nil)
}

func PathValueBool(r *http.Request, key string) Result[bool] {
	value, err := strconv.ParseBool(r.PathValue(key))
	return NewResult(value, err)
}

func JsonDecoder[T any](r *http.Request) Result[*T] {
	var dto T
	err := json.NewDecoder(r.Body).Decode(&dto)
	return NewResult(&dto, err)
}

func PathValueEmail(r *http.Request, key string) Result[string] {
	valeu := r.PathValue(key)
	email, err := mail.ParseAddress(valeu)
	return NewResult(email.Address, err)
}

func PathValueUINT(r *http.Request, key string) Result[uint64] {
	return NewResult(strconv.ParseUint(r.PathValue(key), 10, 64))
}

func PathValueFloat32(r *http.Request, key string) Result[float32] {
	value, err := strconv.ParseFloat(r.PathValue(key), 32)
	return NewResult(float32(value), err)
}

func PathValueFloat64(r *http.Request, key string) Result[float64] {
	return NewResult(strconv.ParseFloat(r.PathValue(key), 64))
}

type MultipartFile struct {
	File   multipart.File
	Header *multipart.FileHeader
}

func GetClaim(r *http.Request, key any) Result[string] {
	value, ok := r.Context().Value(key).(string)
	if !ok {
		return NewResult("", fmt.Errorf("invalid claim %s", key))
	}
	return NewResult(value, nil)
}

func GetClaimUUID(r *http.Request, key any) Result[uuid.UUID] {
	value, ok := r.Context().Value(key).(string)
	if !ok {
		return NewResult(uuid.UUID{}, fmt.Errorf("invalid claim %s", key))
	}
	result, err := uuid.Parse(value)
	if err != nil {
		return NewResult(uuid.UUID{}, fmt.Errorf("invalid claim %s", key))
	}
	return NewResult(result, nil)
}

func FormFile(r *http.Request, key string) Result[*MultipartFile] {
	file, header, err := r.FormFile(key)
	if err != nil {
		return Result[*MultipartFile]{Error: err}
	}
	return Result[*MultipartFile]{
		Value: &MultipartFile{
			File:   file,
			Header: header,
		},
	}
}

func isSliceEmpty(value any) bool {
	if slice, ok := value.([]int); ok {
		return len(slice) == 0
	}
	if slice, ok := value.([]string); ok {
		return len(slice) == 0
	}
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Slice && v.Len() == 0
}

func ReplyJsonOrError(result any, err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		status := http.StatusBadRequest
		if (err.Error() == "record not found") || (err.Error() == "sql: no rows in result set") {
			status = http.StatusNoContent
		}
		http.Error(w, err.Error(), status)
		return
	}
	if result != nil && !isSliceEmpty(result) {
		json.NewEncoder(w).Encode(result)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ReplyTextPlainOrError(result string, err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, result)
}

func ReplyIfError(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func ReplayCreatedOrError(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
