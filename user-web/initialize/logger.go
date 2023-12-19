package initialize

import "go.uber.org/zap"

func InitLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		//"/Users/huberyyang/Desktop/shop-api/user-web/mytest.log",
		"stderr",
	}
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}
