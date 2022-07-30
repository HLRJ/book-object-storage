package locate

import (
	"encoding/json"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//拿到第三个元素的内容，就是<object_name>
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	// 如果没有就返回404
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//json 序列化 并将该地址写入返回信息中
	b, _ := json.Marshal(info)
	w.Write(b)
}
