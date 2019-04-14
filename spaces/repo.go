package spaces

import (
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"
)

// Repository interface specifies database api
type Repository interface {
	Create(u Space) (Space, error)
	GetById(id int) (Space, error)
}

type mysqlSpacesRepository struct {
	DB *sql.DB
}

// NewMysqlSpacesRepository returns a struct that implements the mysqlSpacesRepository interface
func NewMysqlSpacesRepository(db *sql.DB) *mysqlSpacesRepository {
	return &mysqlSpacesRepository{
		DB: db,
	}
}

// Create inserts a new space into the db
func (r *mysqlSpacesRepository) Create(s Space) (Space, error) {
	// TODO: validate & sanitize

	sql, args, err := sq.Insert("spaces").SetMap(sq.Eq{
		"owner_id":    s.OwnerID,
		"title":       s.Title,
		"description": s.Description,
		"address":     s.Address,
		"city":        s.City,
		"province":    s.Province,
		"country":     s.Country,
		"postal_code": s.PostalCode,
	}).ToSql()

	if err != nil {
		log.Printf("error in space repo: %s", err.Error())
		return Space{}, err
	}

	res, err := r.DB.Exec(sql, args...)
	if err != nil {
		log.Printf("error in space repo: %s", err.Error())
		return Space{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("error in space repo: %s", err.Error())
		return Space{}, err
	}

	return Space{ID: uint(id)}, nil
}

// GetById get space by id
func (r *mysqlSpacesRepository) GetById(id int) (Space, error) {
	var space Space

	stmnt, args, err := sq.Select("*").From("Spaces").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		log.Printf("error in spaces repo: %s", err.Error())
		return Space{}, err
	}

	err = r.DB.QueryRow(stmnt, args...).Scan(
		&space.ID,
		&space.OwnerID,
		&space.Title,
		&space.Description,
		&space.Address,
		&space.City,
		&space.Province,
		&space.Country,
		&space.PostalCode,
		&space.CreatedAt,
		&space.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("space %d not found.", space.ID)
			return Space{}, err
		}
		log.Printf("error in space repo: %s", err.Error())
		return Space{}, err
	}

	return space, nil
}
