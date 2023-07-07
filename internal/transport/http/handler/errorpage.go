package handler

import (
	"html/template"
	"net/http"
)

type ErrorBody struct {
	Message string
	Code    int
}

func errorHandler(w http.ResponseWriter, message string, code int) {
	errorBody := ErrorBody{Message: message, Code: code}
	html, err := template.ParseFiles("./ui/error.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	if err = html.Execute(w, errorBody); err != nil {
		// fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}
