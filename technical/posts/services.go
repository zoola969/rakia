package posts

import (
	"errors"
)

var InvalidPostIDError = errors.New("invalid post ID")

type Service interface {
	GetAllPosts() ([]PostRead, error)
	GetPostByID(id int) (PostRead, error)
	CreatePost(req PostCreateUpdate) (PostRead, error)
	UpdatePost(id int, req PostCreateUpdate) (PostRead, error)
	DeletePost(id int) error
}

type PostService struct {
	repo Repository
}

func NewPostService(repo Repository) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (s *PostService) GetAllPosts() ([]PostRead, error) {
	return s.repo.GetAll()
}

func (s *PostService) GetPostByID(id int) (PostRead, error) {
	if id <= 0 {
		return PostRead{}, errors.New("invalid post ID")
	}
	return s.repo.GetByID(id)
}

func (s *PostService) CreatePost(data PostCreateUpdate) (PostRead, error) {
	if err := data.Validate(); err != nil {
		return PostRead{}, err
	}

	return s.repo.Create(data)
}

func (s *PostService) UpdatePost(id int, data PostCreateUpdate) (PostRead, error) {
	if id <= 0 {
		return PostRead{}, InvalidPostIDError
	}

	if err := data.Validate(); err != nil {
		return PostRead{}, err
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		return PostRead{}, err
	}

	return s.repo.Update(id, data)
}

func (s *PostService) DeletePost(id int) error {
	if id <= 0 {
		return errors.New("invalid post ID")
	}
	return s.repo.Delete(id)
}
