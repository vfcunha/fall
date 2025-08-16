package fall

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/mail"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func QueryValue(r *http.Request, key string) Result[string] {
	result := r.URL.Query().Get(key)
	if result == "" {
		return NewResult(result, fmt.Errorf("invalid query value: %s", key))
	}
	return NewResult(r.URL.Query().Get(key), nil)
}

func PathValueUUID(r *http.Request, key string) Result[uuid.UUID] {
	return NewResult(uuid.Parse(r.PathValue(key)))
}

func PathValueInt64(r *http.Request, key string) Result[int64] {
	return NewResult(strconv.ParseInt(r.PathValue(key), 10, 64))
}

func PathValueUint64(r *http.Request, key string) Result[uint64] {
	return NewResult(strconv.ParseUint(r.PathValue(key), 10, 64))
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
	if err == nil {
		var instance any = &dto
		if validator, ok := instance.(validator); ok {
			err = validator.Validate()
		}
	}
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

func Cookie(r *http.Request, name string) Result[*http.Cookie] {
	return NewResult(r.Cookie(name))
}

func BodyTextPlain(r *http.Request) Result[string] {
	responseData, err := io.ReadAll(r.Body)
	return NewResult(string(responseData), err)
}

type MultipartFile struct {
	File   multipart.File
	Header *multipart.FileHeader
}

func GetClaimString(r *http.Request, key any) Result[string] {
	value, ok := r.Context().Value(key).(string)
	if !ok {
		return NewResult("", fmt.Errorf("invalid claim %s", key))
	}
	return NewResult(value, nil)
}

func GetClaim[T any](r *http.Request, key any) Result[T] {
	ctxValue := r.Context().Value(key)
	if ctxValue == nil {
		var zero T
		return NewResult(zero, fmt.Errorf("claim %v not found in context", key))
	}
	value, ok := ctxValue.(T)
	if !ok {
		var zero T
		actualType := fmt.Sprintf("%T", ctxValue)
		expectedType := fmt.Sprintf("%T", zero)
		return NewResult(zero, fmt.Errorf("invalid claim type for %v - expected %s, got %s",
			key, expectedType, actualType))
	}
	return NewResult(value, nil)
}

func GetClaimNumber[T uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64 | float32 | float64 | int](r *http.Request, key any) Result[T] {
	switch v := r.Context().Value(key).(type) {
	case uint8:
		return NewResult(T(v), nil)
	case uint16:
		return NewResult(T(v), nil)
	case uint32:
		return NewResult(T(v), nil)
	case uint64:
		return NewResult(T(v), nil)
	case int8:
		return NewResult(T(v), nil)
	case int16:
		return NewResult(T(v), nil)
	case int32:
		return NewResult(T(v), nil)
	case int64:
		return NewResult(T(v), nil)
	case float32:
		return NewResult(T(v), nil)
	case float64:
		return NewResult(T(v), nil)
	case int:
		return NewResult(T(v), nil)
	default:
		var zero T
		return NewResult(zero, fmt.Errorf("claim %v is not a number (got %T)", key, v))
	}
}

func GetClaimInt64(r *http.Request, key any) Result[int64] {
	value := r.Context().Value(key)
	if value == nil {
		return NewResult(int64(0), fmt.Errorf("invalid claim %s", key))
	}
	return NewResult(value.(int64), nil)
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
		status := http.StatusUnprocessableEntity
		if (err.Error() == "record not found") || (err.Error() == "sql: no rows in result set") {
			status = http.StatusNoContent
		}
		http.Error(w, err.Error(), status)
		return
	}
	if result != nil && !isSliceEmpty(result) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ReplyTextPlainOrError(result string, err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	fmt.Fprint(w, result)
}

func ReplyIfError(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ReplayCreatedOrError(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func NotBlankWithError(value string, err error) Result[string] {
	if len(strings.TrimSpace(value)) == 0 {
		return NewResult(value, err)
	}
	return NewResult(value, nil)
}

func NotBlank(value string) Result[string] {
	return NotBlankWithError(value, fmt.Errorf("string is blank"))
}

func NotBlankEnv(variableName string) Result[string] {
	return NotBlankWithError(os.Getenv(variableName), fmt.Errorf("%s is blank", variableName))
}
