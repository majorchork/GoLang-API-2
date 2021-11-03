package main

import (
	"GoLang-API-2/servers"
	"log"
	"os"
)

func main() {
	// run the server on the port 3000
	port := os.Getenv("PORT")
	if port == ""{
		port = "3000"
	}
	err := servers.Init(port)
	if err != nil{
		log.Fatal("could not start server")
	}

}