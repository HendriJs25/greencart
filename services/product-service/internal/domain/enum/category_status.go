package enum

type CategoryStatus string

const (
	CategoryStatusPublished   CategoryStatus = "published"
	CategoryStatusUnpublished CategoryStatus = "unpublished"
)

func (c CategoryStatus) Bool() bool {
	return c == CategoryStatusPublished
}

func CategoryStatusFromBool(value bool) CategoryStatus {
	if value {
		return CategoryStatusPublished
	}
	return CategoryStatusUnpublished
}
