package main_rabbitmq

import (
	"../../protocol"
	"github.com/streadway/amqp"
	"log"
)

type RabbitClient struct {
	ch *amqp.Channel
	name string
	writequeue amqp.Queue
	readqueue amqp.Queue
	incoming chan protocol.MessageCommand
}

func NewClient() *RabbitClient {
	return &RabbitClient{
		incoming: make(chan protocol.MessageCommand),
	}
}

func (c * RabbitClient) Dial(address string) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	c.ch = ch
	wq, err := ch.QueueDeclare(
		"clientwriting", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	rq, err := ch.QueueDeclare(
		"clientreading"+address, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	c.readqueue = rq
	c.writequeue = wq
	return nil
}

func (c * RabbitClient) Start() {
	msgs, err := c.ch.Consume(
	  c.readqueue.Name, // queue
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
		c.incoming <- protocol.MessageCommand{
			Message: string(d.Body),
		}
	  }
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (c * RabbitClient) SetName(name string) error{
	c.name = name
	return nil
}

func (c * RabbitClient) SendMessage(body string) error {
	var err = c.ch.Publish(
		"",     // exchange
		c.writequeue.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(body+c.name),
		})
	failOnError(err, "Failed to publish a message")
	return nil
}
func (c * RabbitClient) Close(){}
func (c * RabbitClient) Clean(){
	c.incoming = make(chan protocol.MessageCommand)
}
func (c *RabbitClient) Incoming() chan protocol.MessageCommand {
	return c.incoming
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}