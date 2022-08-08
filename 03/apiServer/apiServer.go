package main

import (
	"book-object-storage/03/apiServer/heartbeat"
	"book-object-storage/03/apiServer/locate"
	"book-object-storage/03/apiServer/objects"
	"book-object-storage/03/apiServer/version"

	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", version.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))

}
