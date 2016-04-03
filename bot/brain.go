package bot

import (
	"errors"
	"fmt"
)

var TypeService string = "service"
var TypeJob string = "job"

type Brain struct {
	Stores []StoreRef
}

func NewBrain() Brain {
	return Brain{}
}

func (b *Brain) Lookup(t, load string) Store {
	if b.Stores == nil {
		return nil
	}

	for i := range b.Stores {
		store := b.Stores[i]
		if store.Type == t && store.Load == load {
			return store.Impl
		}
	}

	return nil
}

func (b *Brain) Get(t, load string) (s Store, err error) {
	s = b.Lookup(t, load)
	if s == nil {
		switch t {
		case TypeService:
			service := GetService(load)
			s = service.Store()
			b.Stores = append(b.Stores, StoreRef{Type: t, Load: load, Impl: s})
		default:
			return nil, errors.New(fmt.Sprint("Unknown store type: %s", t))
		}
	}

	return s, nil
}
