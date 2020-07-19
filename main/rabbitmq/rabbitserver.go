package main_rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"sync"
)

type RabbitServer struct {
	ch *amqp.Channel
	writeclients []*amqp.Queue
	readclients amqp.Queue
	initialized bool
	mutex   *sync.Mutex
}

func NewServer() *RabbitServer {
	return &RabbitServer{
		mutex: &sync.Mutex{},
		initialized: false,
	}
}

func (c * RabbitServer) Register(name string) {
	println("registering "+name)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	c.ch = ch
	if !c.initialized {
		q, err := ch.QueueDeclare(
			"clientwriting", // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		failOnError(err, "Failed to declare a queue")
		c.readclients = q
		c.initialized = true
	}
	rq, err := ch.QueueDeclare(
		"clientreading"+name, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	c.mutex.Lock()
	c.writeclients = append(c.writeclients, &rq)
	c.mutex.Unlock()
}

func (c * RabbitServer) Start() {
	msgs, err := c.ch.Consume(
		c.readclients.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			c.broadcast(d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (c * RabbitServer) broadcast(body []byte) {
	for _, q := range c.writeclients {
		c.sendMessage(body, q.Name)
	}

}

func (c * RabbitServer) sendMessage(body []byte, queueName string) {
	var err = c.ch.Publish(
		"",     // exchange
		queueName, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
}