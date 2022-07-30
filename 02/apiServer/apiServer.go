package main

import (
	"book-object-storage/02/apiServer/locate"
	"book-object-storage/02/apiServer/objects"

	"book-object-storage/02/apiServer/heartbeat"
	"log"
	"net/http"
	"os"
)

func main() {
	//启一个协程 心跳监控
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
