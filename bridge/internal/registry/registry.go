package registry

import (
	"math/rand"
	"sync"
	"sync/atomic"

	"bridge/internal/limiter"
)

type Service struct {
	Name        string
	Protocol    string // http | grpc
	Addr        string
	Weight      int
	GrayPercent int
	Healthy     atomic.Bool
	Limiter     *limiter.PerServiceLimiter
}

type Registry struct {
	mu       sync.RWMutex
	services map[string][]*Service
}

func New() *Registry {
	return &Registry{services: make(map[string][]*Service)}
}

func (r *Registry) Register(s *Service, qps int) {
	s.Healthy.Store(true)
	s.Limiter = limiter.NewLimiter(qps, qps)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[s.Name] = append(r.services[s.Name], s)
}

func (r *Registry) Pick(name string) *Service {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := r.services[name]
	if len(list) == 0 {
		return nil
	}

	var healthy []*Service
	for _, s := range list {
		if s.Healthy.Load() {
			healthy = append(healthy, s)
		}
	}
	if len(healthy) == 0 {
		return nil
	}

	var candidates []*Service
	for _, s := range healthy {
		if s.GrayPercent <= 0 || rand.Intn(100) < s.GrayPercent {
			candidates = append(candidates, s)
		}
	}
	if len(candidates) == 0 {
		candidates = healthy
	}

	total := 0
	for _, s := range candidates {
		total += s.Weight
	}
	rnd := rand.Intn(total)
	for _, s := range candidates {
		rnd -= s.Weight
		if rnd < 0 {
			return s
		}
	}
	return candidates[0]
}

func (r *Registry) List(name string) []*Service {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.services[name]
}

func (r *Registry) Iter(f func(*Service)) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, list := range r.services {
		for _, s := range list {
			f(s)
		}
	}
}
