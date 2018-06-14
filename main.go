package main

import (
	"fmt"

	"github.com/eejai42/sassymq-golang/helloworld/demo"
)

func main() {

	world := new(demo.World)
	world.Init("amqps://smqPublic:smqPublic@explore.ssot.me/DEMO")
	world.AddProgrammerHelloHandler(func(actor *demo.ActorBase, payload *demo.Payload) *demo.Payload {
		fmt.Println("Got hello from programmer" + payload.Content)
		return payload
	})

	programmer := new(demo.Programmer)

	programmer.Init("amqps://smqPublic:smqPublic@explore.ssot.me/DEMO")
	payload := programmer.CreatePayload()
	payload.Content = "This is the payload"
	programmer.Hello(payload, func(actor *demo.ActorBase, reply *demo.Payload) *demo.Payload {
		fmt.Println("Got reply from the world")
		return reply
	})

	select {}
}
