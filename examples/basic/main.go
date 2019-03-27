package main

import (
	"log"
	"net/http"

	"github.com/syntaqx/serve"
)

func main() {
	fs := serve.NewFileServer(serve.Options{
		Directory: "../../static",
	})

	log.Print("serve started at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", fs))
}
