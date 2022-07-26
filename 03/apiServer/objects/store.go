package objects

import (
	"io"
	"net/http"
)

func storeObject(r io.Reader, object string) (int, error) {
	//使用流式传输
	stream, err := putStream(object)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	io.Copy(stream, r)
	err = stream.Close()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
