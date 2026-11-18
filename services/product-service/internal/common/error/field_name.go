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

var FieldNameJa = mergeFieldName(CategoryFieldNameJa)

func mergeFieldName(mapsToMerge ...map[string]string) map[string]string {
	result := make(map[string]string)

	for _, m := range mapsToMerge {
		maps.Copy(result, m)
	}
	return result
}
