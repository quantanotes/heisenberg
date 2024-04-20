package worker

type Worker interface {
	Receive(any)
}
