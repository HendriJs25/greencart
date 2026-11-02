package error

import "maps"

var AuthFieldNameJa = map[string]string{
	"Name":                 "名前",
	"Email":                "メールアドレス",
	"Password":             "パスワード",
	"PasswordConfirmation": "パスワード確認",
}

var FieldNameJa = mergeFieldName(AuthFieldNameJa)

func mergeFieldName(mapsToMerge ...map[string]string) map[string]string {
	result := make(map[string]string)

	for _, m := range mapsToMerge {
		maps.Copy(result, m)
	}
	return result
}
