package worker

type Manager struct {
	processes []process
}

func NewManager() *Manager {
	return &Manager{processes: []process{}}
}

func (m *Manager) Add(w Worker) {
	m.processes = append(m.processes, newProcess(w))
}

func (m *Manager) Start() {
	for _, p := range m.processes {
		go p.Start()
	}
}

func (m *Manager) Stop() {
	for _, p := range m.processes {
		p.Stop()
	}
}

func (m *Manager) Send(msg any) {
	for _, p := range m.processes {
		p.Send(msg)
	}
}

func (m *Manager) Size() int {
	return len(m.processes)
}
