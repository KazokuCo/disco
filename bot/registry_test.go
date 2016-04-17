package bot

import (
	"github.com/codegangsta/cli"
	"testing"
)

type testService struct{ Var bool }

func (*testService) Start(Store)          {}
func (*testService) Stop(Store)           {}
func (*testService) Store() Store         { return &struct{}{} }
func (*testService) Command() cli.Command { return cli.Command{} }

func testServiceFactory() Service { return &testService{Var: true} }

func TestRegisterService(t *testing.T) {
	RegisterService("test", testServiceFactory)

	fn, ok := ServiceRegistry["test"]
	if !ok {
		t.Error("RegisterService() doesn't register anything!")
	}

	s, ok := fn().(*testService)
	if !ok {
		t.Error("Wrong type returned")
		return
	}
	if !s.Var {
		t.Error("Service not initialized")
	}
}

func TestGetService(t *testing.T) {
	RegisterService("test", testServiceFactory)

	s, ok := GetService("test").(*testService)
	if !ok {
		t.Error("Wrong type returned")
		return
	}
	if !s.Var {
		t.Error("Service not initialized")
	}
}

type testJob struct{ Var bool }

func testJobFactory() Job { return &testJob{Var: true} }

func TestRegisterJob(t *testing.T) {
	RegisterJob("test", testJobFactory)

	fn, ok := JobRegistry["test"]
	if !ok {
		t.Error("RegisterJob() doesn't register anything!")
	}

	s, ok := fn().(*testJob)
	if !ok {
		t.Error("Wrong type returned")
		return
	}
	if !s.Var {
		t.Error("Job not initialized")
	}
}

func TestGetJob(t *testing.T) {
	RegisterJob("test", testJobFactory)

	s, ok := GetJob("test").(*testJob)
	if !ok {
		t.Error("Wrong type returned")
		return
	}
	if !s.Var {
		t.Error("Job not initialized")
	}
}
