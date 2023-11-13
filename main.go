package main

import (
	server "microservice/server"
)

func main() {
	router := server.Server()
	router.Run()
}
