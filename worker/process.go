package worker

type process struct {
	Meta   Metadata
	worker Worker
	inbox  chan any
}

func newProcess(w Worker, m Metadata) process {
	return process{
		Meta:   m,
		worker: w,
		inbox:  make(chan any, 4096),
	}
}

func (p *process) Start() {
	for msg := range p.inbox {
		p.worker.Receive(msg)
	}
}

func (p *process) Stop() {
	close(p.inbox)
}

func (p *process) Send(msg any) {
	p.inbox <- msg
}
