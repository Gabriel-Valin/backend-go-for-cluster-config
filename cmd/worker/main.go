package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbit, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("rabbitmq: %v", err)
	}
	defer rabbit.Close()

	ch, _ := rabbit.Channel()
	defer ch.Close()
	ch.QueueDeclare("pedidos", true, false, false, false, nil)

	msgs, _ := ch.Consume("pedidos", "", true, false, false, false, nil)
	log.Println("worker esperando mensagens...")
	for msg := range msgs {
		log.Printf("processando pedido: %s", string(msg.Body))
	}
}
