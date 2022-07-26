package heartbeat

import (
	"os"
	"time"
)
import "book-object-storage/src/lib/rabbitmq"

//

func StartHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	// 无限循环中调用rabbitmq.RabbitMQ结构体的Publish方法向apiServers exchange 发送本节点的监听地址。
	for {
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}
