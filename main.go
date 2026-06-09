package main

import (
	"code-judge/internal/service"
	"code-judge/internal/transport/http"
	"code-judge/pkg/logger"
	"context"
	"log/slog"
	"os"
	"runtime"
)

func main() {
	switch os.Getenv("SANDBOX_ROLE") {
	case "child":
		runtime.LockOSThread()

		fd := os.NewFile(3, "sync-pipe")
		fd.Read(make([]byte, 1))
		fd.Close()

		language := os.Args[1]
		filename := os.Args[2]

		service.ExecuteSubmission(language, filename)
	default:
		log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		logger.WithContext(context.Background(), log)

		s := http.NewServer(log)
		log.Info("Server started on port 8080")
		if err := s.ListenAndServe(); err != nil {
			log.Error("error starting server: ", err)
			return
		}
	}
}
