package http

import (
	"code-judge/internal/service"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func NewServer() *http.Server {

	http.Handle("/", http.HandlerFunc(indexHandler))
	http.Handle("/submit", http.HandlerFunc(submitHandler))
	return &http.Server{
		Addr: ":8080",

		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  10 * time.Second,
		WriteTimeout:                 10 * time.Second,
		MaxHeaderBytes:               1 << 20,
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(r.Form["code_submission"][0])
	service.Execute(r.Form["code_submission"][0])
}
