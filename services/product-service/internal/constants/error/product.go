package error

import "errors"

var (
	ErrProductNotFound           = errors.New("product not found")
	ErrAtLeastOneVariantRequired = errors.New("at least one product variant is required")
)
