package spaces

import (
	"github.com/asaskevich/govalidator"
	"github.com/gobuffalo/nulls"
)

// Space models a space
type Space struct {
	ID          uint        `json:"id,omitempty" db:"id"`
	OwnerID     uint        `json:"owner_id,omitempty" db:"owner_id"`
	Title       string      `json:"title,omitempty" db:"title" db:"owner_id" valid:"required~title is required"`
	Description string      `json:"description,omitempty" db:"description" db:"owner_id" valid:"required~description is required"`
	Address     string      `json:"address,omitempty" db:"address" db:"owner_id" valid:"required~address is required"`
	City        string      `json:"city,omitempty" db:"city" db:"owner_id" valid:"required~city is required"`
	Province    string      `json:"province,omitempty" db:"province" db:"owner_id" valid:"required~province is required"`
	Country     string      `json:"country,omitempty" db:"country" db:"owner_id" valid:"required~country is required"`
	PostalCode  string      `json:"postal_code,omitempty" db:"postal_code"`
	Tags        []Tag       `json:"tags,omitempty" db:"-"`
	CreatedAt   *nulls.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   *nulls.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Valid validates a Space struct.
func (s Space) Valid() (bool, error) { return govalidator.ValidateStruct(s) }

// TagCategory is the list of allowed values for a tag's category.
type TagCategory string

const (
	// TagCategoryNone no tag category.
	TagCategoryNone TagCategory = "none"
	// TagCategoryUsage usage tags.
	TagCategoryUsage TagCategory = "usage"
	// TagCategoryDiscipline discipline tags.
	TagCategoryDiscipline TagCategory = "discipline"
)

// Tag models a space tag
type Tag struct {
	ID       uint        `json:"id,omitempty" db:"id" valid:"required~id is required"`
	Title    string      `json:"title,omitempty" db:"title"`
	Category TagCategory `json:"category,omitempty" db:"category"`
}

// Valid validates a Tag struct.
func (s Tag) Valid() (bool, error) { return govalidator.ValidateStruct(s) }
