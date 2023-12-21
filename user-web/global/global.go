package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop-api/user-web/configs"
	"time"
)

var (
	ServerSetting *configs.ServerSettingS
	Trans         ut.Translator
	Location      *time.Location
)
