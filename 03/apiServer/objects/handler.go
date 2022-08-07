package objects

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method // 判断方法的类型，然后执行相应的函数
	if method == http.MethodPut {
		put(w, r)
		return
	}
	if method == http.MethodGet {
		get(w, r)
		return
	}
	if method == http.MethodDelete {
		del(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)

}
