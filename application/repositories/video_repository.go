package repositories

import (
	"encoder/domain"
	"fmt"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type VideoRepository interface {
	Insert(v *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}

type VideoRepositoryDB struct {
	Db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepositoryDB {
	return &VideoRepositoryDB{Db: db}
}

func (repo VideoRepositoryDB) Insert(video *domain.Video) (*domain.Video, error) {

	id := uuid.NewV4().String()

	if video.ID == "" {
		video.ID = id
	}

	err := repo.Db.Create(video).Error

	if err != nil {
		return nil, err
	}

	return video, nil
}

func (r VideoRepositoryDB) Find(id string) (*domain.Video, error) {
	var video domain.Video

	r.Db.Preload("Jobs").First(&video, "id = ?", id)

	if video.ID == "" {
		return nil, fmt.Errorf("video does not exists")
	}

	return &video, nil
}
