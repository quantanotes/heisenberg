package worker

type Manager struct {
	processes []process
}

type SubManager struct {
	processes []process
}

func NewManager() *Manager {
	return &Manager{processes: []process{}}
}

func (m *Manager) Add(w Worker, meta Metadata) {
	m.processes = append(m.processes, newProcess(w, meta))
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

func (m *Manager) WhereJobEq(job int) *SubManager {
	processes := []process{}
	for _, p := range m.processes {
		if p.Meta.Job == job {
			processes = append(processes, p)
		}
	}
	return &SubManager{processes}
}

func (sm *SubManager) Send(msg any) {
	for _, p := range sm.processes {
		p.Send(msg)
	}
}
