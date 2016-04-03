package bot

import (
	"errors"
	"fmt"
)

type Job interface {
}

type JobRef struct {
	Load string
	Impl Job
}

func (ref *JobRef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var loadRef struct{ Load string }

	if err := unmarshal(&loadRef); err != nil {
		return err
	}

	ref.Load = loadRef.Load
	ref.Impl = GetJob(loadRef.Load)

	if ref.Impl == nil {
		return errors.New(fmt.Sprintf("Job not found: %s", loadRef.Load))
	}

	return nil
}
