package repository

import (
	"errors"

	"github.com/Sarthak-D97/go_stuAPI/entity"
	"gorm.io/gorm"
)

type Repository interface {
	Create(student entity.Student) (entity.Student, error)
	GetByID(id int64) (*entity.Student, error)
	List() ([]entity.Student, error)
	Update(id int64, student entity.Student) error
	Delete(id int64) error
}

type gormRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(student entity.Student) (entity.Student, error) {
	if err := r.db.Create(&student).Error; err != nil {
		return entity.Student{}, err
	}
	return student, nil
}

func (r *gormRepository) GetByID(id int64) (*entity.Student, error) {
	var student entity.Student
	if err := r.db.First(&student, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &student, nil
}

func (r *gormRepository) List() ([]entity.Student, error) {
	var students []entity.Student
	if err := r.db.Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

func (r *gormRepository) Update(id int64, student entity.Student) error {
	student.ID = int(id)
	return r.db.Save(&student).Error
}

func (r *gormRepository) Delete(id int64) error {
	return r.db.Delete(&entity.Student{}, id).Error
}
