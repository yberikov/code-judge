package main

import (
	"code-judge/internal/transport/http"
	"fmt"
)

func main() {
	s := http.NewServer()
	fmt.Println("Server started")
	if err := s.ListenAndServe(); err != nil {
		fmt.Println(err)
		return
	}
}
