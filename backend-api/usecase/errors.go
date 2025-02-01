package usecase

import "errors"

var (
	ErrInvalidEmail       = errors.New("メールアドレスが無効です")
	ErrInvalidCredentials = errors.New("メールアドレスまたはパスワードが正しくありません")
	ErrEmailExists        = errors.New("このメールアドレスは既に登録されています")
	ErrHashPassword       = errors.New("パスワードのハッシュ化に失敗しました")
	ErrCreateUser         = errors.New("ユーザーの作成に失敗しました")
	ErrGenerateToken      = errors.New("トークンの生成に失敗しました")
)
