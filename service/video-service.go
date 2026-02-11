package service

import (
	"github.com/Sarthak-D97/go_stuAPI/entity"
	"github.com/Sarthak-D97/go_stuAPI/repository"
)

type VideoService interface {
	Save(entity.Video) entity.Video
	Update(video entity.Video)
	Delete(video entity.Video)
	FindAll() ([]entity.Video, error)
}

type videoService struct {
	videoRepository repository.VideoRepository
}

func NewVideoService(repo repository.VideoRepository) VideoService {
	return &videoService{
		videoRepository: repo,
	}
}
func (s *videoService) Update(video entity.Video) {
	s.videoRepository.Update(video)
}
func (s *videoService) Delete(video entity.Video) {
	s.videoRepository.Delete(video)
}

func (s *videoService) Save(video entity.Video) entity.Video {
	s.videoRepository.Save(video)
	return video
}
func (s *videoService) FindAll() ([]entity.Video, error) {
	return s.videoRepository.FindAll()
}
