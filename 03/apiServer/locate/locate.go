package locate

import (
	"book-object-storage/src/lib/rabbitmq"
	"os"
	"strconv"
	"time"
)

func Locate(name string) string {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	// 消息队列发布消息
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		// 一秒钟后关闭消息队列，相当于是一个超时机制，避免无休止的等待
		time.Sleep(time.Second)
		q.Close()
	}()
	// 消息内容传给masg
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

// Exist 检查Locate是否为空，来判断对象是否存在。
func Exist(name string) bool {
	return Locate(name) != ""
}
