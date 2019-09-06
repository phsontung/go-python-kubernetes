package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	DateTime time.Time `json:"datetime"`
	Message  string    `json:"message"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string) // broadcast channel

var ampqCh *amqp.Channel
var amqpQ amqp.Queue

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		// msg := <-broadcast

		msgs, err := ampqCh.Consume(
			amqpQ.Name, // queue
			"",         // consumer
			true,       // auto ack
			true,       // exclusive
			false,      // no local
			false,      // no wait
			nil,        // args
		)
		failOnError(err, "Failed to register a consumer")

		for d := range msgs {
			log.Printf(" [x] %s - %s", d.RoutingKey, d.Body)

			response := Message{
				DateTime: time.Now(),
				Message:  string(d.Body),
			}
			// res2B, _ := json.Marshal(response)

			// log.Printf("Response  %s\n", msg)
			// log.Printf("Response message %s\n", res2B)
			// Send it out to every client that is currently connected
			for client := range clients {
				err := client.WriteJSON(response)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func setupConsumer() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	// defer ch.Close()

	err = ch.ExchangeDeclare(
		"topic_logs", // name
		"topic",      // type
		false,        // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,       // queue name
		"#",          // routing key
		"topic_logs", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	amqpQ = q
	ampqCh = ch
}

func main() {
	http.HandleFunc("/ws/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		clients[conn] = true

		defer conn.Close()

		for {
			// Read message from browser
			_, msg, err := conn.ReadMessage()
			if err != nil {
				delete(clients, conn)
				fmt.Printf("Connection closed.\n")
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			// if err = conn.WriteMessage(msgType, msg); err != nil {
			// 	return
			// }

			// Send the newly received message to the broadcast channel
			// broadcast <- string(msg)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/src/notification/internal/templates/websocket.html")
	})

	setupConsumer()
	// Start listening for incoming chat messages
	go handleMessages()

	fmt.Printf("Start server at :8081\n")
	http.ListenAndServe(":8081", nil)
}
