package objects

import (
	"book-object-storage/02/apiServer/heartbeat"
	"book-object-storage/src/lib/objectstream"
	"fmt"
)

// 用对象的名字，找到一个随机地可用的节点地址，然后将其
func putStream(object string) (*objectstream.PutStream, error) {
	// 找到一个能连接的节点，随机就可以
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	return objectstream.NewPutStream(server, object), nil
}
