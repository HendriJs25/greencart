package error

import "errors"

var (
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrParentCategoryMissing = errors.New("parent category does not exist")
)
