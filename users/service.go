package users

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/LiamPimlott/spaces/lib"
)

type UsersService interface {
	Create(u User) (User, error)
	GetById(id int) (User, error)
	// Login(u User) (Token, error)
}

type usersService struct {
	repo   UsersRepository
	secret string
}

// NewUsersService will return a struct that implements the UsersService interface
func NewUsersService(repo UsersRepository, secret string) *usersService {
	return &usersService{
		repo:   repo,
		secret: secret,
	}
}

// Create creates a new user and issues a token
func (s *usersService) Create(u User) (User, error) {

	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		log.Printf("error in user service: %s", err.Error())
		return User{}, err
	}

	u.Password = string(pass)

	usr, err := s.repo.Create(u)
	if err != nil {
		return User{}, err
	}

	token, err := utils.GenerateToken(usr.Username, s.secret)
	if err != nil {
		return User{}, err
	}

	usr.Token = token

	return usr, nil
}

// GetById retrieves a user by their id
func (s *usersService) GetById(id int) (User, error) {
	usr, err := s.repo.GetById(id)
	if err != nil {
		return User{}, err
	}
	return usr, nil
}

// TODO: Login validates an email/pass and return a token
