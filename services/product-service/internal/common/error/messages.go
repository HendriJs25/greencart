package error

import (
	"errors"
	errConstant "product-service/internal/constants/error"
)

var ErrValidator = map[string]string{
	"required": "%sは必須です",
	"notblank": "%sは空白のみでは登録できません",
}

func MessageJa(err error) string {
	switch {
	case errors.Is(err, errConstant.ErrBadRequest):
		return "入力内容に誤りがあります"
	case errors.Is(err, errConstant.ErrInvalidCSRFToken):
		return "CSRFトークンの検証に失敗しました"
	case errors.Is(err, errConstant.ErrInvalidToken):
		return "認証に失敗しました"
	case errors.Is(err, errConstant.ErrFileTooBig):
		return "ファイルサイズが大きすぎます"
	case errors.Is(err, errConstant.ErrUnsupportedFileExtension):
		return "対応していないファイル形式です"
	case errors.Is(err, errConstant.ErrForbidden):
		return "この操作を実行する権限がありません"
	default:
		return "サーバー内部でエラーが発生しました"
	}
}
