package error

import "maps"

var AuthFieldNameJa = map[string]string{
	"Email":    "メールアドレス",
	"Password": "パスワード",
}

var FieldNameJa = mergeFieldName(AuthFieldNameJa)

func mergeFieldName(mapsToMerge ...map[string]string) map[string]string {
	result := make(map[string]string)

	for _, m := range mapsToMerge {
		maps.Copy(result, m)
	}
	return result
}
