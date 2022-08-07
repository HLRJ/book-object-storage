package objects

import (
	"book-object-storage/src/lib/es/es8"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	// 计算请求头的哈希值
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 将r.body 存入
	code, err := storeObject(r.Body, url.PathEscape(hash))
	if err != nil {
		log.Println(err)
		w.WriteHeader(code)
		return
	}
	if code != http.StatusOK {
		w.WriteHeader(code)
		return
	}
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size := utils.GetSizeFromHeader(r.Header)
	err = es8.AddVersion(name, hash, size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
