package fall

import "gorm.io/gorm"

type Entity[K any] interface {
	GetID() K
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

func (r *GormRepository[E, K]) Delete(entity *K) error {
	return r.DB.Delete(entity).Error
}

func (r *GormRepository[E, K]) Update(entity *E) error {
	err := r.DB.Save(entity).Error
	if err != nil {
		return err
	}
	return nil
}
