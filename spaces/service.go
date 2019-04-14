package spaces

import (
	"log"
)

// Service interface to spaces service
type Service interface {
	Create(u Space) (Space, error)
	GetByID(id int) (Space, error)
}

type spacesService struct {
	repo Repository
}

// NewSpacesService will return a struct that implements the spacesService interface
func NewSpacesService(repo Repository) *spacesService {
	return &spacesService{
		repo: repo,
	}
}

// Create creates a new Space and issues a token
func (s *spacesService) Create(u Space) (Space, error) {
	// TODO santitize and validate input
	space, err := s.repo.Create(u)
	if err != nil {
		log.Printf("error creating space: %s\n", err)
		return Space{}, err
	}
	return space, nil
}

// GetByID retrieves a space by their id
func (s *spacesService) GetByID(id int) (Space, error) {
	space, err := s.repo.GetById(id)
	if err != nil {
		log.Printf("error getting space: %s\n", err)
		return Space{}, err
	}
	return space, nil
}
