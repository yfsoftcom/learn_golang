package runner

type Runner struct {
	tasks []func(int)
}

func New() *Runner {
	return &Runner{}
}

func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

func (r *Runner) Start() {

	for i, task := range r.tasks {
		task(i)
	}
}
