package controllers

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

type logger struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func (log logger) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (log *logger) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (log *logger) notFound(w http.ResponseWriter) {
	log.clientError(w, http.StatusNotFound)
}
