package error

import "maps"

var AuthFieldNameJa = map[string]string{
	"Name":                 "名前",
	"Email":                "メールアドレス",
	"Password":             "パスワード",
	"PasswordNew":          "新しいパスワード",
	"PasswordConfirmation": "パスワード確認",
}

var UserFieldNameJa = map[string]string{
	"Name":            "名前",
	"Phone":           "電話番号",
	"Lat":             "緯度",
	"Lng":             "経度",
	"Address":         "住所",
	"CurrentPassword": "現在のパスワード",
}

var CustomerFieldNameJa = map[string]string{
	"IsVerified": "認証状態",
}

var FieldNameJa = mergeFieldName(AuthFieldNameJa, UserFieldNameJa, CustomerFieldNameJa)

func mergeFieldName(mapsToMerge ...map[string]string) map[string]string {
	result := make(map[string]string)

	for _, m := range mapsToMerge {
		maps.Copy(result, m)
	}
	return result
}
