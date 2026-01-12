package capabilities

import (
	"fmt"
	"github.com/hyperjumptech/grule-rule-engine/ast"
)

type Registry struct {
	capabilities map[string]Capability
}

func NewRegistry() *Registry {
	return &Registry{
		capabilities: make(map[string]Capability),
	}
}

func (r *Registry) Register(cap Capability) error {
	if _, exists := r.capabilities[cap.Name()]; exists {
		return fmt.Errorf("capability %s already registered", cap.Name())
	}
	r.capabilities[cap.Name()] = cap
	return nil
}

func (r *Registry) Get(name string) Capability {
	return r.capabilities[name]
}

func (r *Registry) BuildDataContext(imei string) (ast.IDataContext, error) {
	dc := ast.NewDataContext()
	for _, cap := range r.capabilities {
		if err := cap.Initialize(imei); err != nil {
			return nil, err
		}
		dc.Add(cap.GetDataContextName(), cap)
	}
	return dc, nil
}