package demo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sassymq/sassymq-golang-helloworld/demo/errors"
	"github.com/streadway/amqp"
)

type ActorBase struct {
	actor_name     string
	sender_id      uuid.UUID
	connection     *amqp.Connection
	channel        *amqp.Channel
	listener_queue *amqp.Queue
	reply_queue    *amqp.Queue
	handlers       map[string]interface{}
	reply_handlers map[string]interface{}
	queue          string
	exchange       string
	connStr        string
	ActorBaseActions
}

type ActorBaseActions interface {
	Init()
	CreateListener()
}

func (p *ActorBase) CreatePayload() *Payload {
	payload := new(Payload)
	pid, err := uuid.NewUUID()
	errors.DefaultWithPanic(err, "Failed to create uuid")
	payload.PayloadId = pid.String()
	payload.SenderId = p.sender_id.String()
	return payload
}

func (p *ActorBase) CloseAll() {
	err := p.channel.Close()
	errors.DefaultWithPanic(err, "Failed to close actorbase channel")
	err = p.connection.Close()
	errors.DefaultWithPanic(err, "Failed to close actorbase connection")
}

func (p *ActorBase) Init(actor string, connStr string) {
	p.actor_name = actor
	p.connStr = connStr
	p.queue = actor + ".all"
	p.exchange = actor + "mic"
	senderId, err := uuid.NewUUID()
	errors.DefaultWithPanic(err, "Failed to create sender id")
	p.sender_id = senderId
	err = p.SetConnection(connStr)
	errors.DefaultWithPanic(err, "Failed to initiate actorbase connection")
	err = p.SetChannel()
	errors.DefaultWithPanic(err, "Failed to initiate actorbase channel")
	err = p.SetQos()
	errors.DefaultWithPanic(err, "Issue setting channel Qos")
	err = p.CreateListener()
	errors.DefaultWithPanic(err, "Failed to initiate queue listener")
	err = p.ResetHandlers()
	errors.DefaultWithPanic(err, "Failed to initiate handlers map")
}

func (p *ActorBase) SetConnection(address string) error {
	conn, err := amqp.Dial(address)
	if err == nil {
		p.connection = conn
	}
	return err
}

func (p *ActorBase) SetChannel() error {
	ch, err := p.connection.Channel()
	if err == nil {
		p.channel = ch
	}
	return err
}

func (p *ActorBase) SetQos() error {
	err := p.channel.Qos(
		1,     // Prefetch count
		0,     // Prefetch size
		false, // Global
	)
	return err
}

func (p *ActorBase) CreateListener() error {
	q, err := p.channel.QueueDeclare(
		p.queue, // Name
		true,    // Durable
		false,   // Delete when unused
		false,   // Exclusive
		false,   // No-wait
		nil,     // Arguments
	)
	if err == nil {
		p.listener_queue = &q
		go p.listenForever()
		q, err = p.channel.QueueDeclare(
			"",    // Name
			true,  // Durable
			false, // Delete when unused
			false, // Exclusive
			false, // No-wait
			nil,   // Arguments
		)
		if err == nil {
			p.reply_queue = &q
			go p.awaitRepliesForever()
		}

	}
	return err
}

func (p *ActorBase) listenForever() {
	msgs, err := p.channel.Consume(
		p.listener_queue.Name, // Queue
		"",    // Consumer
		false, // Auto-ack
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,   // Args
	)
	errors.Default(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			fmt.Printf("%s  received a message: %s", p.queue, d.RoutingKey)
			if value, ok := p.handlers[string(d.RoutingKey)]; ok {
				var payload *Payload
				json.Unmarshal(d.Body, &payload)
				replyPayload := value.(func(*ActorBase, *Payload) *Payload)(p, payload)
				if replyPayload != nil {
					replyPayloadJson, err := json.Marshal(replyPayload)
					errors.DefaultWithPanic(err, "Can't marshal reply payload")
					p.channel.Publish("", d.ReplyTo, false, false,
						amqp.Publishing{
							DeliveryMode: amqp.Persistent,
							ContentType:  "text/plain",
							Body:         []byte(replyPayloadJson),
						})
				}
			} else {
				errors.Default(nil, "Handler not defined")
			}
			d.Ack(false)
		}
	}()

	forever := make(chan bool)
	<-forever
}

func (p *ActorBase) awaitRepliesForever() {
	msgs, err := p.channel.Consume(
		p.reply_queue.Name, // Queue
		"",                 // Consumer
		false,              // Auto-ack
		false,              // Exclusive
		false,              // No-local
		false,              // No-wait
		nil,                // Args
	)
	errors.Default(err, "Failed to start consuming reply queue")

	go func() {
		for d := range msgs {
			var payload *Payload
			err := json.Unmarshal(d.Body, &payload)
			errors.DefaultWithPanic(err, "Can't unmarshal payload to json")
			fmt.Printf("%s received a REPLY message: %s", p.queue, payload.ReplyContent)
			if replyHandler, ok := p.reply_handlers[payload.PayloadId]; ok {
				delete(p.reply_handlers, payload.PayloadId)
				replyHandler.(func(*ActorBase, *Payload) *Payload)(p, payload)
			}

			// Wait for messages
			d.Ack(false)
		}
	}()

	forever := make(chan bool)
	<-forever
}

func (p *ActorBase) Publish(exchange string, key string, mandatory bool, immediate bool, publishing amqp.Publishing) error {
	err := p.channel.Publish(
		exchange,  // exchange
		key,       // routing key
		mandatory, // mandatory
		immediate, // immediate
		publishing)
	return err
}

func (p *ActorBase) QueueDeclare(name string, durable bool, should_delete bool, exclusive bool, no_wait bool, args amqp.Table) (amqp.Queue, error) {
	q, err := p.channel.QueueDeclare(
		name,          // Name
		durable,       // Durable
		should_delete, // Delete when unused
		exclusive,     // Exclusive
		no_wait,       // No-wait
		args,          // Arguments
	)
	return q, err
}

func (p *ActorBase) AddHandler(key string, value interface{}) error {
	p.handlers[key] = value
	return nil
}

func (p *ActorBase) RemoveHandler(key string) error {
	delete(p.handlers, key)
	return nil
}

func (p *ActorBase) ResetHandlers() error {
	p.handlers = make(map[string]interface{})
	p.reply_handlers = make(map[string]interface{})
	return nil
}

func (p *ActorBase) SendPayload(routingKey string, payload *Payload, replyHandler interface{}) error {
	payloadId, err := uuid.NewUUID()
	errors.DefaultWithPanic(err, "Failed to create uuid")
	payload.PayloadId = payloadId.String()
	payloadJson, err := json.Marshal(payload)
	errors.DefaultWithPanic(err, "Can't marshal payload as Json")

	err = p.channel.Publish(
		p.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode:  amqp.Persistent,
			ContentType:   "text/plain",
			CorrelationId: payload.PayloadId,
			ReplyTo:       p.reply_queue.Name,
			Body:          []byte(payloadJson),
		})

	if replyHandler != nil {
		p.reply_handlers[payload.PayloadId] = replyHandler
		go p.WaitForReply(payload, 10)
	}

	return err
}

func (actor *ActorBase) WaitForReply(payload *Payload, timeout_seconds int) {
	time.Sleep(time.Duration(timeout_seconds) * time.Second)
	if replyPayload, ok := actor.reply_handlers[payload.PayloadId]; ok {
		if replyPayload != nil {
			fmt.Println("TIMED OUT WAITING FOR REPLY" + payload.PayloadId)
			delete(actor.reply_handlers, payload.PayloadId)
		}
	}
}
