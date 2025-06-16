package providers

import "sync"

// Registry holds registered providers.
var Registry = struct {
	sync.RWMutex
	providers map[string]Provider
}{providers: map[string]Provider{}}

// Register adds a provider to the registry.
func Register(p Provider) {
	Registry.Lock()
	defer Registry.Unlock()
	Registry.providers[p.Name()] = p
}

// Get returns a provider by name.
func Get(name string) Provider {
	Registry.RLock()
	defer Registry.RUnlock()
	return Registry.providers[name]
}
