package bot

type ServiceFactory func() Service
type JobFactory func() Job

var ServiceRegistry map[string]ServiceFactory
var JobRegistry map[string]JobFactory

func init() {
	ServiceRegistry = make(map[string]ServiceFactory)
	JobRegistry = make(map[string]JobFactory)
}

func RegisterService(id string, fn ServiceFactory) {
	ServiceRegistry[id] = fn
}

func RegisterJob(id string, fn JobFactory) {
	JobRegistry[id] = fn
}

func GetService(id string) Service {
	fn := ServiceRegistry[id]
	if fn != nil {
		return fn()
	}
	return nil
}

func GetJob(id string) Job {
	fn := JobRegistry[id]
	if fn != nil {
		return fn()
	}
	return nil
}
