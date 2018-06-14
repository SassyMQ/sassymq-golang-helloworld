package demo

type Programmer struct {
	ActorBase
	ProgrammerActions
}

type ProgrammerActions interface {
	Init()
	
    Hello()
}

func (p *Programmer) Init(connStr string) {
	p.ActorBase.Init("programmer", connStr)
}


func (p *Programmer) Hello(payload *Payload, replyHandler interface{}) error {
	err := p.SendPayload("world.general.programmer.hello", payload, replyHandler)
	return err
}



func (p *Programmer) AddWorldWassupHandler(value interface{}) error {
	err := p.AddHandler("programmer.general.world.wassup", value)
	return err
}

                