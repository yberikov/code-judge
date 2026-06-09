package http

import (
	"code-judge/internal/service"
	"code-judge/internal/transport/http/middleware"
	"code-judge/internal/transport/http/model"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

func NewServer(logger *slog.Logger) *http.Server {

	http.Handle("/", middleware.Chain(
		http.HandlerFunc(indexHandler),
		middleware.LoggerMiddleware(logger),
	))

	http.Handle("/submit", middleware.Chain(
		http.HandlerFunc(submitHandler),
		middleware.LoggerMiddleware(logger),
	))

	return &http.Server{
		Addr:                         ":8080",
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
	var req model.SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := service.InitSubmission(r.Context(), req.CodeSubmission)
	if err != nil {
		return
	}
	response, err := json.Marshal(result)
	w.Write(response)
}
