package main_rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"net"
)

type RabbitClient struct {
	ch chan struct{}
	writequeue amqp.Queue
	readqueue amqp.Queue
	incoming chan struct{}
}

func (c * RabbitClient) Dial(name string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	c.ch = ch

	wq, err := ch.QueueDeclare(
		"clientwrinting", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	rq, err := ch.QueueDeclare(
		"clientreading"+name, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	c.readqueue = rq
	c.writequeue = wq
}

func (c * RabbitClient) start() {
	msgs, err := ch.Consume(
	  q.Name, // queue
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
		c.incoming <- d.Body
	  }
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (c * RabbitClient) sendMessage(body string) {
	err = c.ch.Publish(
		"",     // exchange
		c.wq.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}



func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

