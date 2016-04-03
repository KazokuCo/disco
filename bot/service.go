package bot

import (
	"errors"
	"fmt"
)

type Service interface {
	// Creates a blank Store. Loading stores is handled for you, but need to be
	// marshalable into YAML (using: https://gopkg.in/yaml.v2).
	Store() Store

	// Starts the service.
	Start(store Store)

	// Interactively authorizes the bot with the remote service.
	// Return whether changes to the bot's brain should be persisted.
	Login(store Store) bool
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
