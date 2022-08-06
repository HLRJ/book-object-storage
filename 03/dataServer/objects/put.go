package objects

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// 将r中body复制到本地创建的文件中
func put(w http.ResponseWriter, r *http.Request) {
	// 创建文件 前面的为目录
	f, err := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + strings.Split(r.URL.EscapedPath(), "/")[2])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, r.Body)
}
