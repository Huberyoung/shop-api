package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"os"
	"strings"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"shop-api/user-web/custom_validator"
	"shop-api/user-web/global"
	"shop-api/user-web/initialize"
)

func main() {

	// 通过目录判断是否是在家
	dir, _ := os.Getwd()
	viper.Set("is_home", true)
	if strings.Contains(dir, "CloudDocs") {
		viper.Set("is_home", false)
	}

	// 初始化配置
	if err := initialize.InitConfig(); err != nil {
		zap.S().Fatalf("initialize.InitConfig err:%s\n", err)
	}

	host, port := global.ServerSetting.HttpHost, global.ServerSetting.HttpPort
	zap.S().Infof("启动端口号%s:%s\n", host, port)

	// 初始化日志
	initialize.InitLogger()

	// 初始化错误翻译
	if err := initialize.InitTrans("zh"); err != nil {
		zap.S().Fatalf("initialize.InitTrans err:%s\n", err)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		registerFn := func(ut ut.Translator) error {
			return ut.Add("mobile", "{0}手机号码非法!", true)
		}
		translationFn := func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		}

		_ = v.RegisterValidation("mobile", custom_validator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, registerFn, translationFn)
	}

	// 初始化路由
	router := initialize.InitRouter()
	if err := router.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		zap.S().Fatalf("router.Run err:%s\n", err)
	}
}
