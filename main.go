package main

import "wilin/src/server"

func main() {
	router := server.New()
	router.Run("0.0.0.0:8080")
}
