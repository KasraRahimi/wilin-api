package main

import "wilin/src/server"

func main() {
	router := server.New()
	router.Run("localhost:8080")
}
