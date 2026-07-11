package fall

import (
	"math"

	"gorm.io/gorm"
)

type Entity[K any] interface {
	GetID() K
}

type Page[T any] struct {
	Content          []T   `json:"content"`
	Number           int   `json:"number"`
	Size             int   `json:"size"`
	TotalElements    int64 `json:"totalElements"`
	TotalPages       int   `json:"totalPages"`
	First            bool  `json:"first"`
	Last             bool  `json:"last"`
	NumberOfElements int   `json:"numberOfElements"`
}

func NewPage[T any](content []T, number, size int, totalElements int64) Page[T] {
	totalPages := 0
	if size > 0 {
		totalPages = int(math.Ceil(float64(totalElements) / float64(size)))
	}

	numberOfElements := len(content)

	return Page[T]{
		Content:          content,
		Number:           number,
		Size:             size,
		TotalElements:    totalElements,
		TotalPages:       totalPages,
		First:            number == 0,
		Last:             number >= totalPages-1 && totalPages > 0,
		NumberOfElements: numberOfElements,
	}
}

type GormRepository[E any, K any] struct {
	DB *gorm.DB `fall:"database"`
}

func (r *GormRepository[E, K]) Save(entity *E) (*E, error) {
	err := r.DB.Create(entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *GormRepository[E, K]) SaveAll(list []*E) ([]*E, error) {
	err := r.DB.Create(list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *GormRepository[E, K]) FindAll() (*[]E, error) {
	var entity *[]E
	err := r.DB.Find(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *GormRepository[E, K]) FindById(id K) (*E, error) {
	var entity *E
	err := r.DB.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *GormRepository[E, K]) FindByIds(ids []K) ([]*E, error) {
	var entities []*E
	err := r.DB.Find(&entities, ids).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *GormRepository[E, K]) Delete(entity *E) error {
	return r.DB.Delete(entity).Error
}

func (r *GormRepository[E, K]) DeleteById(id K) error {
	err := r.DB.Delete(new(E), id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormRepository[E, K]) Update(entity *E) (*E, error) {
	err := r.DB.Save(entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *GormRepository[E, K]) FindPage(query *gorm.DB, page, size int) (Page[E], error) {
	if page < 0 {
		page = 0
	}
	if size < 1 {
		size = 10
	}

	var total int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return Page[E]{}, err
	}

	var results []E
	offset := page * size
	if err := query.Offset(offset).Limit(size).Find(&results).Error; err != nil {
		return Page[E]{}, err
	}

	return NewPage(results, page, size, total), nil
}
