package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/streadway/amqp"
)

type clientDetails struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

//Client hearbeat... keeps pushing its name and ip address
func main() {

	fmt.Printf("broadcasting IP address via MQ: %s\n", "rabitMQ")
	ip := GetOutboundIP()
	fmt.Printf("broadcasting IP address via MQ: %s\n", ip)

	//connect to rabbitmq
	conn, err := amqp.Dial("amqp://alchemy_apache:Password1@crdc-001uatcbe1:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"alchemy-cluster", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare a queue")

	for {
		client := clientDetails{Name: GetClientName(), Address: ip.String()}
		body, err := json.Marshal(client)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(string(body)),
			})
		fmt.Printf("broadcasting IP address via MQ: %s\n", ip.String())
		failOnError(err, "Failed to publish a message")
		time.Sleep(30 * time.Second)
	}

	log.Printf("Starting consumer")

	//consume message

	failOnError(err, "Failed to register a consumer")
	log.Printf("Displaying messages")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func GetClientName() string {
	count := len(os.Args[1:])
	var arg string
	if count == 0 {
		arg = "clientUnkown"
	} else {
		arg = os.Args[1]
		if arg != "" {
			arg = arg
		}
	}
	return arg
}
