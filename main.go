package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {

	host := flag.String("host", ":8080", "host del servidor")
	flag.Parse()

	hub := newHub() //

	go hub.run() // lleva registro del lo goroutines

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	// recodar que cada conexion se manda llamar como goroutine
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		handleConnection(hub, writer, request)
	})

	fmt.Println("Server on " + *host)

	log.Fatal(http.ListenAndServe(*host, nil))

}
