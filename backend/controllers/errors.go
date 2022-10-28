package controllers

// type logger struct {
// 	ErrorLog *log.Logger
// 	InfoLog  *log.Logger
// }

// func (l logger) serverError(w http.ResponseWriter, err error) {
// 	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
// 	l.ErrorLog.Output(2, trace)

// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// }

// func (l *logger) clientError(w http.ResponseWriter, status int) {
// 	http.Error(w, http.StatusText(status), status)
// }

// func (l *logger) notFound(w http.ResponseWriter) {
// 	l.clientError(w, http.StatusNotFound)
// }
