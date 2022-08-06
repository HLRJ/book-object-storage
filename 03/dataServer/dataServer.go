package main

import (
	"book-object-storage/03/dataServer/heartbeat"
	"book-object-storage/03/dataServer/locate"
	"book-object-storage/03/dataServer/objects"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
