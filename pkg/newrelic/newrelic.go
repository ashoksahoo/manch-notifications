package newrelic

import "github.com/newrelic/go-agent"

var App newrelic.Application

func init() {
	config := newrelic.NewConfig("Manch Notification Service", "eu01xx038a0925d9d0e0945a787a04d06fd5c0a9")
	App, _ = newrelic.NewApplication(config)
}
