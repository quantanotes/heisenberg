package heisenberg

type Heisenberg struct {
}

func New() Heisenberg {
	return Heisenberg{}
}

func (h *Heisenberg) Run() {
	h.repl()
}

func (h *Heisenberg) open(path string) {}