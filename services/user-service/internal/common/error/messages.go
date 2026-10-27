package error

import (
	"errors"
	errConstant "user-service/internal/constants/error"
)

var ErrValidator = map[string]string{
	"required":      "%sは必須です",
	"email":         "%sの形式は正しくありません",
	"max":           "%sは%s文字以内で入力してください",
	"min":           "%sは%s文字以上で入力してください",
	"validpassword": "%sは8文字以上で、大文字・小文字・数字をそれぞれ1文字以上含めてください",
	"eqfield":       "%sが%sと一致していません",
	"notblank":      "%sは空白のみでは登録できません",
}

func MessageJa(err error) string {
	switch {
	case errors.Is(err, errConstant.ErrBadRequest):
		return "入力内容に誤りがあります"
	case errors.Is(err, errConstant.ErrInvalidEmailOrPassword):
		return "メールアドレスまたはパスワードが正しくありません"
	case errors.Is(err, errConstant.ErrInvalidCSRFToken):
		return "CSRFトークンの検証に失敗しました"
	case errors.Is(err, errConstant.ErrInvalidToken):
		return "認証に失敗しました"
	default:
		return "サーバー内部でエラーが発生しました"
	}
}
