package heartbeat

import (
	"book-object-storage/src/lib/rabbitmq"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// 用于接收和处理来自数据服务节点的心跳消息
// dataServers 在整个包内可见，用于缓存所有的数据服务节点
var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

// ListenHeartbeat 建一个消息队列，将其绑定到apiServers 通过go channel 监听每一个来自数据服务节点的心跳消息，
//将该消息的正文内容：数据服务节点的监听地址作为    map的键，使用锁来保证并发安全（dataServers）
//收到消息的时间作为值存入dataServers
func ListenHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("apiServers")
	c := q.Consume()
	go removeExpiredDataServer()
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}
func removeExpiredDataServer() {
	// 每隔5s扫描一遍dataServers，并清除其中超过10s没收到心跳消息的数据服务节点。
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			//如果加上10s在当前时间之前，相当于10秒内没响应，直接删除
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

// GetDataServers 遍历dataServers 并返回当前所有的数据服务节点。
func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	// 地址切片
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}

// 会在当前所有的数据服务节点中随机选出一个节点并返回，如果当前数据服务节点为空，则返回空字符串
func ChooseRandomDataServer() string {
	ds := GetDataServers()
	n := len(ds)
	if n == 0 {
		return ""
	}
	return ds[rand.Intn(n)]
}
