package bot

import (
	"errors"
	"fmt"
)

type Service interface {
	Start()
}

type ServiceRef struct {
	Load string
	Impl Service
}

func (ref *ServiceRef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var loadRef struct{ Load string }

	if err := unmarshal(&loadRef); err != nil {
		return err
	}

	ref.Load = loadRef.Load
	ref.Impl = GetService(loadRef.Load)

	if ref.Impl == nil {
		return errors.New(fmt.Sprintf("Service not found: %s", loadRef.Load))
	}

	return unmarshal(ref.Impl)
}
