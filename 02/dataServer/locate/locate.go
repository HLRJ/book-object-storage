package locate

import (
	"book-object-storage/src/lib/rabbitmq"
	"os"
	"strconv"
)

// Locate 如果存在目标文件，返回true，不存在返回false
func Locate(name string) bool {
	_, err := os.Stat(name)
	//如果err的值为nil，说明文件或文件夹存在
	//返回的错误类型 使用 os.IsNotExist() 判断为true，说明文件或文件夹不存在
	return !os.IsNotExist(err)
}

func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	//c为一个channel
	c := q.Consume()
	for msg := range c {
		//将msg.Body的字符串去除引号（json带双引号），消息队列里面传的是apiserver接口服务发送过来需要定位的对象名字
		object, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		//locate里面的是完整的地址
		if Locate(os.Getenv("STORAGE_ROOT") + "/objects/" + object) {
			//如果存在，则调用Send方法向消息的发送方返回本服务节点的  监听地址，表示该对象存在于本服务节点上。
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}
