package users

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/LiamPimlott/spaces/lib"
)

type UsersService interface {
	Create(u User) (User, error)
	Login(u User) (User, error)
	GetById(id int) (User, error)
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
	// TODO santitize and validate input
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return User{}, err
	}

	u.Password = string(pass)

	usr, err := s.repo.Create(u)
	if err != nil {
		log.Printf("error creating user: %s\n", err)
		return User{}, err
	}

	token, err := utils.GenerateToken(usr.Username, s.secret)
	if err != nil {
		log.Printf("error creating user: %s\n", err)
		return User{}, err
	}

	usr.Token = token

	return usr, nil
}

// Login validates an email/pass and return a token
func (s *usersService) Login(u User) (User, error) {
	usr, err := s.repo.GetPassword(u.Email)
	if err != nil {
		log.Printf("error logging in user: %s\n", err)
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(u.Password))
	if err != nil {
		log.Printf("error logging in user: %s\n", err)
		return User{}, err
	}

	token, err := utils.GenerateToken(usr.Username, s.secret)
	if err != nil {
		log.Printf("error logging in user: %s\n", err)
		return User{}, err
	}

	return User{Token: token}, nil
}

// GetById retrieves a user by their id
func (s *usersService) GetById(id int) (User, error) {
	usr, err := s.repo.GetById(id)
	if err != nil {
		log.Printf("error getting user: %s\n", err)
		return User{}, err
	}
	return usr, nil
}
