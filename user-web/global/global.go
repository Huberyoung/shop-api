package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop-api/user-web/configs"
)

var (
	ServerSetting *configs.ServerSettingS
	Trans         ut.Translator
)
