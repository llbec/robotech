package router

import (
	"net/url"
	"sync"
)

type Service struct {
	Name       string
	URL        *url.URL
	HealthURL  string
	MetricsURL string
}

type Registry struct {
	mu       sync.RWMutex
	services map[string]*Service
}

func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]*Service),
	}
}

// 添加或更新服务
func (r *Registry) AddService(name, rawURL, health, metrics string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[name] = &Service{
		Name: name, URL: parsed, HealthURL: health, MetricsURL: metrics,
	}
	return nil
}

func (r *Registry) GetService(name string) (*Service, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.services[name]
	return s, ok
}

func (r *Registry) ListServices() []*Service {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := []*Service{}
	for _, s := range r.services {
		list = append(list, s)
	}
	return list
}
