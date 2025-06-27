package global

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var Log = logrus.New()

func ConfigureLogger() {
	// 设置日志级别
	Log.SetLevel(logrus.InfoLevel)

	// 创建日志文件
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.SetOutput(file)
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}

	// 设置JSON格式
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
}
