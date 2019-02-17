package users

import (
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"
)

type UsersRepository interface {
	Create(u User) (User, error)
	GetPassword(email string) (User, error)
	GetById(id int) (User, error)
}

type mysqlUsersRepository struct {
	DB *sql.DB
}

// NewMysqlUsersRepository returns a struct that implements the UsersRepository interface
func NewMysqlUsersRepository(db *sql.DB) *mysqlUsersRepository {
	return &mysqlUsersRepository{
		DB: db,
	}
}

func (r *mysqlUsersRepository) Create(u User) (User, error) {
	// TODO: validate & sanitize

	sql, args, err := sq.Insert("users").SetMap(sq.Eq{
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"username":   u.Username,
		"email":      u.Email,
		"password":   u.Password,
	}).ToSql()

	if err != nil {
		log.Printf("error in user repo: %s", err.Error())
		return User{}, err
	}

	res, err := r.DB.Exec(sql, args...)
	if err != nil {
		log.Printf("error in user repo: %s", err.Error())
		return User{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("error in user repo: %s", err.Error())
		return User{}, err
	}

	return User{ID: uint(id)}, nil
}

func (r *mysqlUsersRepository) GetPassword(email string) (User, error) {
	var usr User

	sql, args, err := sq.Select("password").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()

	if err != nil {
		log.Printf("error in user repo: %s", err.Error())
		return User{}, err
	}

	err = r.DB.QueryRow(sql, args...).Scan(&usr.Password)
	if err != nil {
		log.Printf("error in user repo: %s", err.Error())
		return User{}, err
	}

	return usr, nil
}

func (r *mysqlUsersRepository) GetById(id int) (User, error) {
	var usr User

	err := r.DB.QueryRow("select id, email from users where id = ?", id).Scan(&usr.ID, &usr.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("user %d not found.", usr.ID)
			return User{}, err
		}
		log.Printf("error in user repo: %s", err.Error())
		return User{}, err
	}

	log.Printf("user %d retrieved.", usr.ID)
	return usr, nil
}
