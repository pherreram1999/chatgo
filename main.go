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

	fs := http.FileServer(http.Dir("./dist"))

	http.Handle("/", fs)

	// recodar que cada conexion se manda llamar como goroutine
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		handleConnection(hub, writer, request)
	})

	fmt.Println("Server on " + *host)

	log.Fatal(http.ListenAndServe(*host, nil))

}
