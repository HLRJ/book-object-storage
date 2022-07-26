package rabbitmq

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

//封装自己的rabbitmq 简化接口

type RabbitMQ struct {
	channel  *amqp.Channel
	conn     *amqp.Connection
	Name     string
	exchange string
}

//创建一个RabbitMQ结构体

func New(s string) *RabbitMQ {
	conn, e := amqp.Dial(s)
	if e != nil {
		panic(e)
	}

	ch, e := conn.Channel()
	if e != nil {
		panic(e)
	}

	q, e := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if e != nil {
		panic(e)
	}

	mq := new(RabbitMQ)
	mq.channel = ch
	mq.conn = conn
	mq.Name = q.Name
	return mq
}

//将消息队列和一个exchange绑定，所有发往exchange的消息都能在自己的消息队列中被接收到。

func (q *RabbitMQ) Bind(exchange string) {
	e := q.channel.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		exchange, // exchange
		false,
		nil)
	if e != nil {
		panic(e)
	}
	q.exchange = exchange
}

//send 可以往某个消息队列发送消息

func (q *RabbitMQ) Send(queue string, body interface{}) {
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}
	e = q.channel.Publish("",
		queue,
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		panic(e)
	}
}

//Publish 方法可以往某个exchange发送消息

func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}
	e = q.channel.Publish(exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		panic(e)
	}
}

//Consume 方法用于生成一个接受消息的 go channel
//使客户程序可以通过Go语言的原生机制接收队列中的消息

func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
	c, e := q.channel.Consume(q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if e != nil {
		panic(e)
	}
	return c
}

//Close 方法用于关闭消息队列
func (q *RabbitMQ) Close() {
	q.channel.Close()
	q.conn.Close()
}
