package endpoint

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetupLogger() *logrus.Logger {
	logger := logrus.New()

	file, err := os.OpenFile("api.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Fatal("Failed to write to file, using default stderr")
	}
	logger.SetOutput(file)

	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}

func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		ctx.Next()

		latency := time.Since(startTime)

		statusCode := ctx.Writer.Status()

		logger.WithFields(logrus.Fields{
			"method":     ctx.Request.Method,
			"path":       ctx.Request.URL.Path,
			"status":     statusCode,
			"latency":    latency,
			"client_ip":  ctx.ClientIP(),
			"user_agent": ctx.Request.UserAgent(),
		}).Info("Api Endpoint")
	}
}
