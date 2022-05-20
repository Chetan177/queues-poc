package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

var (
	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchangeName = flag.String("exchange", "test-topic-exchange", "Durable AMQP exchange name")
	exchangeType = flag.String("exchange-type", "topic", "Exchange type - direct|fanout|topic|x-custom")
	routingKey   = flag.String("key", "test-key", "AMQP routing key")
	body         = flag.String("body", "foobar", "Body of message")
	reliable     = flag.Bool("reliable", true, "Wait for the publisher confirmation before exiting")
)

func init() {
	flag.Parse()
}

func main() {
	if err := publish(*uri, *exchangeName, *exchangeType, *routingKey, *body, *reliable); err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("published %dB OK", len(*body))
}

func publish(amqpURI, exchange, exchangeType, routingKey, body string, reliable bool) error {

	// This function dials, connects, declares, publishes, and tears down,
	// all in one go. In a real service, you probably want to maintain a
	// long-lived connection as state, and publish against that.

	log.Printf("dialing %q", amqpURI)

	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}
	defer connection.Close()

	log.Printf("got Connection, getting Channel")

	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring %q Exchange (%q)", exchangeType, exchange)
	if err := channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	// if reliable {
	// 	log.Printf("enabling publishing confirms.")
	// 	if err := channel.Confirm(false); err != nil {
	// 		return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
	// 	}

	// 	confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	// 	defer confirmOne(confirms)
	// }

	log.Printf("declared Exchange, publishing %dB body (%q)", len(body), body)

	pri := []int{5,1,9}
	count := 0 
	for i := 0; ; i++ {
		routingKey := "sales.call"
		priority := pri[count]
		count++
		if count == 3 {
			count = 0
		}
		expire := "60000"
		// if i%2 == 0 {
		// 	expire = "1000"
		// }
		log.Printf("publishing %dB body (%q), counter %d, expire %s, priority %d", len(body), body, i, expire, priority)
		if err = channel.Publish(
			exchange,   // publish to an exchange
			routingKey, // routing to 0 or more queues
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				Headers:         amqp.Table{},
				ContentType:     "text/plain",
				ContentEncoding: "",
				Body:            []byte(routingKey + " " + body + " " + fmt.Sprintf("priority=%d", priority) + " " + "expire="+ expire),
				DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
				Priority:        uint8(priority), // 0-9
				Expiration:      expire,
				// a bunch of application/implementation-specific fields
			},
		); err != nil {
			return fmt.Errorf("Exchange Publish: %s", err)
		}
		if i%3 == 0{
			// time.Sleep(time.Second * 10)
		}

		time.Sleep(time.Second * 1)
	}

}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
