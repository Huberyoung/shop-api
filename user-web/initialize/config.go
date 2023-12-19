package initialize

import (
	"shop-api/user-web/configs"
	"shop-api/user-web/global"
)

func InitConfig() error {
	setting, err := configs.NewSetting()
	if err != nil {
		return err
	}

	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	return nil
}
