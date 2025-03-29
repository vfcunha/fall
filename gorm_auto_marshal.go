package fall

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type GormAutoMarshal[T any] struct{}

func (j *GormAutoMarshal[T]) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSON value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

func (j GormAutoMarshal[T]) Value() (driver.Value, error) {
	return json.Marshal(j)
}
