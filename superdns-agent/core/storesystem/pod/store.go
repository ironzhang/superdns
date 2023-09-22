package pod

import "sync"

type Store struct {
	mu   sync.RWMutex
	apps map[string]*Application
}

func NewStore() *Store {
	return &Store{apps: make(map[string]*Application)}
}

func (p *Store) SetApplication(app *Application) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.apps[app.name] = app
}

func (p *Store) GetApplication(name string) (*Application, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	app, ok := p.apps[name]
	return app, ok
}
