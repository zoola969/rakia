package posts

import (
	"github.com/go-playground/validator/v10"
)

type PostRead struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type PostCreateUpdate struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Author  string `json:"author" validate:"required"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (d *PostCreateUpdate) Validate() error {
	return validate.Struct(d)
}
