package enum

type CategoryStatus string

const (
	CategoryStatusPublished   CategoryStatus = "有効"
	CategoryStatusUnpublished CategoryStatus = "無効"
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
