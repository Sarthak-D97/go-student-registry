package repository

import (
	"github.com/Sarthak-D97/go_stuAPI/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VideoRepository interface {
	Save(video entity.Video) (entity.Video, error)
	Update(video entity.Video) error
	Delete(video entity.Video) error
	FindAll() ([]entity.Video, error)
	CloseDB() error
}

type database struct {
	connection *gorm.DB
}

func NewVideoRepository() VideoRepository {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to start database Connection")
	}
	err = db.AutoMigrate(&entity.Video{}, &entity.Person{})
	if err != nil {
		panic("failed to migrate database")
	}

	return &database{
		connection: db,
	}
}

func (db *database) CloseDB() error {
	sqlDB, err := db.connection.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (db *database) Save(video entity.Video) (entity.Video, error) {
	err := db.connection.Create(&video).Error
	return video, err
}

func (db *database) Update(video entity.Video) error {
	return db.connection.Save(&video).Error
}

func (db *database) Delete(video entity.Video) error {
	return db.connection.Delete(&video).Error
}

func (db *database) FindAll() ([]entity.Video, error) {
	var videos []entity.Video
	err := db.connection.Preload(clause.Associations).Find(&videos).Error
	return videos, err
}
