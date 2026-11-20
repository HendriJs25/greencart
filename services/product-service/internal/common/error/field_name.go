package error

import "maps"

var CategoryFieldNameJa = map[string]string{
	"ParentSlug":  "親カテゴリースラッグ",
	"Name":        "カテゴリー名",
	"Icon":        "アイコン",
	"Status":      "ステータス",
	"Slug":        "スラッグ",
	"Description": "説明",
}

var ProductFieldNameJa = map[string]string{
	"ProductName":        "商品名",
	"CategorySlug":       "カテゴリー",
	"Unit":               "単位",
	"ProductDescription": "商品説明",
	"Status":             "ステータス",
	"VariantDetail":      "商品バリエーション",
	"Stock":              "在庫数",
	"Weight":             "重量",
	"SalePrice":          "販売価格",
	"RegularPrice":       "通常価格",
}

var FieldNameJa = mergeFieldName(CategoryFieldNameJa, ProductFieldNameJa)

func mergeFieldName(mapsToMerge ...map[string]string) map[string]string {
	result := make(map[string]string)

	for _, m := range mapsToMerge {
		maps.Copy(result, m)
	}
	return result
}
