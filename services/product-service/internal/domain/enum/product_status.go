package enum

type ProductStatus string

var (
	ProductStatusDraft    ProductStatus = "下書き"
	ProductStatusActive   ProductStatus = "有効"
	ProductStatusInactive ProductStatus = "無効"
)

var (
	ProductStatusDraftDB    = "draft"
	ProductStatusInactiveDB = "inactive"
	ProductStatusActiveDB   = "active"
)

func (p ProductStatus) ConvertToDBValue() *string {
	switch p {
	case "有効":
		return &ProductStatusActiveDB
	case "無効":
		return &ProductStatusInactiveDB
	default:
		return &ProductStatusDraftDB
	}
}

func ConvertFromDBValue(value string) ProductStatus {
	switch value {
	case "active":
		return ProductStatusActive
	case "inactive":
		return ProductStatusInactive
	default:
		return ProductStatusDraft
	}
}
