package main

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"github.com/armon/go-metrics/prometheus"

	"github.com/deciphernow/gm-fabric-go/metrics/sinkobserver"
	"github.com/deciphernow/gm-fabric-go/metrics/subject"
)

func getPrometheusObserverIfNeeded(logger zerolog.Logger) ([]subject.Observer, error) {
	if !viper.GetBool("report_statsd") {
		return nil, nil
	}

	prometheusSink, err := prometheus.NewPrometheusSink()
	if err != nil {
		return nil, errors.Wrap(err, "prometheus.NewPrometheusSink")
	}

	sinkObserver := sinkobserver.New(
		prometheusSink,
		viper.GetDuration("prometheus_mem_interval"),
	)

	logger.Debug().Str("service", "{{.ServiceName}}").
		Dur("interval", viper.GetDuration("prometheus_mem_interval")).
		Msg("reporting Prometheus")

	return []subject.Observer{sinkObserver}, nil
}
