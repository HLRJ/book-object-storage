package version

import (
	"book-object-storage/src/lib/es/es8"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	// 只有get方法才能继续执行下去
	if method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from := 0                                          //从第一开始查
	size := 1000                                       //一次查1000个
	name := strings.Split(r.URL.EscapedPath(), "/")[2] //拿到文件的名字
	for {
		// 拿到一个数组存到metas
		metas, err := es8.SearchAllVersions(name, from, size)
		//出错结束
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// 遍历metas  里面存的该对象所有的版本
		for i := range metas {
			body, _ := json.Marshal(metas[i])
			w.Write(body)
			w.Write([]byte("\n"))
		}
		// 遍历结束
		if len(metas) != size {
			return
		}
		// 更新起始值
		from += size
	}
}
