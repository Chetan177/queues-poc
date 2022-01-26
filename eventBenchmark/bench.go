package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fiorix/go-eventsocket/eventsocket"
)

func main() {
	c, err := eventsocket.Dial("localhost:8021", "ClueCon")
	if err != nil {
		log.Fatal(err)
	}

	counter := 0
	c.Send("events plain ALL")

	timeout := time.After(30 * time.Second)
	for {

		select {
		case <-timeout:
			timeout = time.After(30 * time.Second)
			fmt.Println("Time ", timeout, " Counter :", counter)
			counter = 0
		default:
		}

		_, err := c.ReadEvent()
		if err == nil {
			counter++
		}

	}

}
