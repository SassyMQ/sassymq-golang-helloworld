package main

import (
	"fmt"

	"github.com/sassymq/sassymq-golang-helloworld/demo"
)

func main() {

	world := new(demo.World)
	world.Init("amqp://guest:guest@localhost/DEMO")
	world.AddProgrammerHelloHandler(func(actor *demo.ActorBase, payload *demo.Payload) *demo.Payload {
		fmt.Println("" + payload.Content)
		fmt.Println("Got hello from programmer" + payload.Content)
		return payload
	})

	programmer := new(demo.Programmer)

	programmer.Init("amqp://guest:guest@localhost/DEMO")
	payload := programmer.CreatePayload()
	payload.Content = "This is the payload"
	programmer.Hello(payload, func(actor *demo.ActorBase, reply *demo.Payload) *demo.Payload {
		fmt.Println("")
		fmt.Println("Got reply from the world")
		return reply
	})

	select {}
}
