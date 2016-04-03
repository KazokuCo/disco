package bot

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
)

type Store interface {
}

type StoreRef struct {
	Type string
	Load string
	Impl Store
}

func (ref *StoreRef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var loadRef struct {
		Type string
		Load string
	}

	if err := unmarshal(&loadRef); err != nil {
		return err
	}

	ref.Type = loadRef.Type
	ref.Load = loadRef.Load

	switch ref.Type {
	case "service":
		service := GetService(ref.Load)
		if service == nil {
			return errors.New(fmt.Sprintf("Unknown service: %s", ref.Load))
		}
		ref.Impl = service.Store()
		return unmarshal(ref.Impl)
	default:
		return errors.New(fmt.Sprintf("Unknown store type: %s", ref.Type))
	}
}

func (ref StoreRef) MarshalYAML() (interface{}, error) {
	// This is lazy and terrible. Basically just flattening everything into a map by converting
	// to/from YAML, then sticking type info onto it.
	implData, err := yaml.Marshal(ref.Impl)
	if err != nil {
		return nil, err
	}
	implMap := make(map[interface{}]interface{})
	if err = yaml.Unmarshal(implData, implMap); err != nil {
		return nil, err
	}
	implMap["type"] = ref.Type
	implMap["load"] = ref.Load
	return implMap, nil
}
