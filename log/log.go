package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type CustomResponseWriter struct {
	responseWriter http.ResponseWriter
	StatusCode     int
}

func NewCustomResponseWriter(w http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{w, http.StatusOK}
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	return w.responseWriter.Write(b)
}

func (w *CustomResponseWriter) Header() http.Header {
	return w.responseWriter.Header()
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.responseWriter.WriteHeader(statusCode)
}

func SetupLogger(fileName string) *logrus.Logger {
	logger := logrus.New()

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Fatal("Failed to write to file, using default stderr")
	}
	logger.SetOutput(file)

	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func LoggerMiddleware(logger *logrus.Logger, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		crw := NewCustomResponseWriter(w)
		handler.ServeHTTP(crw, r)
		latency := time.Since(startTime)

		logger.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     crw.StatusCode,
			"latency":    latency,
			"client_ip":  ReadUserIP(r),
			"user_agent": r.UserAgent(),
			"access_key": r.Header.Get("Authorization"),
		}).Info("Api Endpoint")
	}
}
