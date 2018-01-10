package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	gometrics "github.com/armon/go-metrics"

	"github.com/deciphernow/gm-fabric-go/metrics/sinkobserver"
	"github.com/deciphernow/gm-fabric-go/metrics/subject"
)

func getStatsdObserverIfNeeded(logger zerolog.Logger) ([]subject.Observer, error) {
	if !viper.GetBool("report_statsd") {
		return nil, nil
	}

	statsdSink, err := gometrics.NewStatsiteSink(
		fmt.Sprintf(
			"%s:%d",
			viper.GetString("statsd_server_host"),
			viper.GetInt("statsd_server_port"),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "gometrics.NewStatsiteSink")
	}

	sinkObserver := sinkobserver.New(
		statsdSink,
		viper.GetDuration("statsd_mem_interval"),
	)

	logger.Debug().Str("service", "{{.ServiceName}}").
		Str("host", viper.GetString("statsd_server_host")).
		Int("port", viper.GetInt("statsd_server_port")).
		Dur("interval", viper.GetDuration("statsd_mem_interval")).
		Msg("reporting statsd")

	return []subject.Observer{sinkObserver}, nil
}
