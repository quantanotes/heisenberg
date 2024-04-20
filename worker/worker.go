package worker

type Worker interface {
	Receive(any)
}

type Metadata struct {
	Job int
}
