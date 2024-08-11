package repository

import "gorm.io/gorm"

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) GetAll(db *gorm.DB, entity *T) error {
	return db.Find(entity).Error
}

func (r *Repository[T]) FindById(db *gorm.DB, entity *T, id any) error {
	return db.Find(entity, id).Error
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T, id any) error {
	return db.Model(entity).Where("id = ?", id).Updates(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}
