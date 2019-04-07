package fastjob

type config struct {
	registry *JobRegistry
	logger   Logger
}

func NewConfig(registry *JobRegistry) *config {
	if registry == nil {
		panic("registry cannot be nil")
	}
	return &config{
		registry: registry,
		logger:   &standardLogger{},
	}
}

func (c *config) WithLogger(logger Logger) *config {
	if logger == nil {
		panic("logger cannot be nil")
	}
	c.logger = logger
	return c
}
