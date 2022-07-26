package main

import (
	"book-object-storage/02/object"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.StarHearbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", object.Handler)
	// 由于是本地运行，将ip端口改成本地ip
	log.Fatal(http.ListenAndServe("127.0.0.1"+os.Getenv("LISTEN_ADDRESS"), nil))
}
