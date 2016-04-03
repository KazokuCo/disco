package bot

import (
	"errors"
	"fmt"
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
		return unmarshal(&ref.Impl)
	default:
		return errors.New(fmt.Sprintf("Unknown store type: %s", ref.Type))
	}
}
