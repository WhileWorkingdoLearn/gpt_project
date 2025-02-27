package main

import (
	"net/http"

	"github.com/WhileCodingDoLearn/gpt_project/internal/auth"
	"github.com/WhileCodingDoLearn/gpt_project/internal/response"
	logger "github.com/WhileCodingDoLearn/gpt_project/log"
	"github.com/sirupsen/logrus"
)

type CustomSeverMux struct {
	smux *http.ServeMux
	log  *logrus.Logger
}

func CustomSmux(logFilePath string) *CustomSeverMux {
	return &CustomSeverMux{
		smux: http.NewServeMux(),
		log:  logger.SetupLogger(logFilePath),
	}
}

func (cmux *CustomSeverMux) Handle(pattern string, handler http.Handler) {
	cmux.smux.Handle(pattern, handler)
}

func (cmux *CustomSeverMux) HandleFunc(pattern string, handler func(w http.ResponseWriter, req *http.Request)) {

	cmux.smux.HandleFunc(pattern, logger.LoggerMiddleware(cmux.log, handler))
}
func (cmux *CustomSeverMux) Handler(r *http.Request) (h http.Handler, pattern string) {
	return cmux.smux.Handler(r)
}

func (cmux *CustomSeverMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cmux.smux.ServeHTTP(w, r)
}

func Authorized(keyType string, key string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetToken(r.Header, keyType)
		if err != nil {
			response.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if len(key) == 0 {
			response.ResponseWithoutBody(w, http.StatusInternalServerError)
			return
		}
		if token != key {
			response.ResponseWithoutBody(w, http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	}
}
