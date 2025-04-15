package main

import (
	"log"
	"wilin/src/server"
)

func main() {
	router := server.New()
	err := router.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
}
