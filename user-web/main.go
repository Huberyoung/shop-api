package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"os"
	"shop-api/user-web/custom_validator"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"

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
		_ = v.RegisterValidation("mobile", custom_validator.ValidateMobile)
	}

	// 初始化路由
	router := initialize.InitRouter()
	if err := router.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		zap.S().Fatalf("router.Run err:%s\n", err)
	}
}
