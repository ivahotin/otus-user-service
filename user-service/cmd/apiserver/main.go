package main

import (
	"log"

	"example.com/arch/user-service/internal/apiserver"
)

func main() {
	server := apiserver.NewServer()
	err := server.Start()
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
