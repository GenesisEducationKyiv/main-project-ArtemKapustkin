package main

import "rabbitmq-log"

func main() {
	consumer := logs_consumer.NewRabbitMQConsumer()
	consumer.LogBindingMessages()
}
