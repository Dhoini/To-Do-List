package main

import (
	"fmt"
	"net/http"
)

func main() {
	App := App()

	server := http.Server{
		Addr:    ":8080",
		Handler: App,
	}

	fmt.Println("Listening on port 8080")

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
