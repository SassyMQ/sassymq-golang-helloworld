package demo

type World struct {
	ActorBase
	WorldActions
}

type WorldActions interface {
	Init()
	
    Wassup()
}

func (p *World) Init(connStr string) {
	p.ActorBase.Init("world", connStr)
}


func (p *World) Wassup(payload *Payload, replyHandler interface{}) error {
	err := p.SendPayload("programmer.general.world.wassup", payload, replyHandler)
	return err
}



func (p *World) AddProgrammerHelloHandler(value interface{}) error {
	err := p.AddHandler("world.general.programmer.hello", value)
	return err
}

                