package readymanager

import "sync"

type ReadyManager struct {
	ready     bool
	readyCh   chan struct{}
	readyChMu sync.Mutex
}

func (r *ReadyManager) ReadyCh() chan struct{} {
	r.readyChMu.Lock()
	defer r.readyChMu.Unlock()
	if r.readyCh == nil {
		r.readyCh = make(chan struct{})
	}
	return r.readyCh
}

func (r *ReadyManager) Ready() bool {
	return r.ready
}

func (r *ReadyManager) SetReady() {
	r.ReadyCh()

	if !r.ready {
		r.readyChMu.Lock()
		defer r.readyChMu.Unlock()
		if !r.ready {
			r.ready = true
			close(r.readyCh)
		}
	}

}
