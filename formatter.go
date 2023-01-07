package yamlfmt

import "fmt"

type Formatter interface {
	Type() string
	Format(yamlContent []byte) ([]byte, error)
}

type Factory interface {
	Type() string
	NewFormatter(config map[string]interface{}) (Formatter, error)
}

type Registry struct {
	registry    map[string]Factory
	defaultType string
}

func NewFormatterRegistry(defaultFactory Factory) *Registry {
	return &Registry{
		registry: map[string]Factory{
			defaultFactory.Type(): defaultFactory,
		},
		defaultType: defaultFactory.Type(),
	}
}

func (r *Registry) Add(f Factory) {
	r.registry[f.Type()] = f
}

func (r *Registry) GetFactory(fType string) (Factory, error) {
	factory, ok := r.registry[fType]
	if !ok {
		return nil, fmt.Errorf("no formatter registered with type \"%s\"", fType)
	}
	return factory, nil
}

func (r *Registry) GetDefaultFactory() (Factory, error) {
	factory, ok := r.registry[r.defaultType]
	if !ok {
		return nil, fmt.Errorf("no default formatter registered for type \"%s\"", r.defaultType)
	}
	return factory, nil
}
