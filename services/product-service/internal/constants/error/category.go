package error

import "errors"

var (
	ErrCategoryAlreadyExists = errors.New("category already exists")
)
