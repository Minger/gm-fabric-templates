package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"google.golang.org/grpc"

	gometrics "github.com/armon/go-metrics"

	"github.com/deciphernow/gm-fabric-go/metrics/gmfabricsink"
	"github.com/deciphernow/gm-fabric-go/metrics/gometricsobserver"
	"github.com/deciphernow/gm-fabric-go/metrics/grpcmetrics"
	"github.com/deciphernow/gm-fabric-go/metrics/grpcobserver"
	ms "github.com/deciphernow/gm-fabric-go/metrics/metricsserver"
	"github.com/deciphernow/gm-fabric-go/metrics/subject"

	"{{.ConfigPackage}}"
	"{{.MethodsPackage}}"

	// we don't use this directly, but need it in vendor for gateway grpc plugin
	_ "github.com/golang/glog"
	_ "github.com/grpc-ecosystem/grpc-gateway/runtime"
)

func main() {
	var tlsMetricsConf *tls.Config
	var tlsServerConf *tls.Config
	var err error
	var zkCancels []zkCancelFunc

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stderr})

	logger.Info().Str("service", "{{.ServiceName}}").Msg("starting")

	ctx, cancelFunc := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		for _, f := range zkCancels {
			f()
		}
	}()
	
	logger.Debug().Str("service", "{{.ServiceName}}").Msg("initializing config")
	if err = {{.ConfigPackageName}}.Initialize(); err != nil {
		logger.Fatal().AnErr("{{.ConfigPackageName}}.Initialize()", err).Msg("")
	}

	if tlsMetricsConf, err = buildMetricsTLSConfigIfNeeded(logger); err != nil {
		logger.Fatal().AnErr("buildMetricsTLSConfigIfNeeded", err).Msg("")
	}

	if tlsServerConf, err = buildServerTLSConfigIfNeeded(logger); err != nil {
		logger.Fatal().AnErr("buildServerTLSConfigIfNeeded", err).Msg("")
	}

	ctx = putOauthInCtxIfNeeded(ctx)

	logger.Debug().Str("service", "{{.ServiceName}}").
		Str("host", viper.GetString("grpc_server_host")).
		Int("port", viper.GetInt("grpc_server_port")).
		Msg("creating listener")

	lis, err := net.Listen(
		"tcp",
		fmt.Sprintf(
			"%s:%d",
			viper.GetString("grpc_server_host"),
			viper.GetInt("grpc_server_port"),
		),
	)
	if err != nil {
		logger.Fatal().AnErr("net.Listen", err).Msg("")
	}

	grpcObserver := grpcobserver.New(viper.GetInt("metrics_cache_size"))
	goMetObserver := gometricsobserver.New()
	observers := []subject.Observer{grpcObserver, goMetObserver}

	statsdObserver, err := getStatsdObserverIfNeeded(logger)
	if err != nil {
		logger.Fatal().AnErr("getStatsdObserverIfNeeded", err).Msg("")
	}
	observers = append(observers, statsdObserver...)
	
	prometheusObserver, err := getPrometheusObserverIfNeeded(logger)
	if err != nil {
		logger.Fatal().AnErr("getPrometheusObserverIfNeeded", err).Msg("")
	}
	observers = append(observers, prometheusObserver...)
	
	logger.Debug().Str("service", "{{.ServiceName}}").
		Str("host", viper.GetString("metrics_server_host")).
		Int("port", viper.GetInt("metrics_server_port")).
		Msg("starting metrics server")

	mux := http.NewServeMux()
	mux.Handle(
		viper.GetString("metrics_dashboard_uri_path"), 
		ms.NewDashboardHandler(grpcObserver.Report, goMetObserver.Report),
	)
	if viper.GetBool("report_prometheus") {
		mux.Handle(
			viper.GetString("metrics_prometheus_uri_path"), 
			ms.NewPrometheusHandler(),
		)
	}
	
	mServer := ms.NewMetricsServer(
		fmt.Sprintf(
			"%s:%d",
			viper.GetString("metrics_server_host"),
			viper.GetInt("metrics_server_port"),
		),		
		tlsMetricsConf,
	)
	mServer.Handler = mux
	
	if mServer.TLSConfig == nil {
		go mServer.ListenAndServe()
	} else {
		go mServer.ListenAndServeTLS("", "")
	}
	
	cancels, err := notifyZkOfMetricsIfNeeded(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("zk metrics announcement")
	}

	zkCancels = append(
		zkCancels,
		cancels...,
	)

	metricsChan := subject.New(ctx, observers...)

	sink := gmfabricsink.New(metricsChan)
	gometrics.NewGlobal(gometrics.DefaultConfig("{{.ServiceName}}"), sink)

	hostName, err := os.Hostname()
	if err != nil {
		logger.Error().Err(err).Msg("unable to determin host name: program continues")
	}
	statsTags := []string{
		subject.JoinTag("service", "{{.ServiceName}}"),
		subject.JoinTag("host", hostName),
	}
	opts := []grpc.ServerOption{
		grpc.StatsHandler(
			grpcmetrics.NewStatsHandlerWithTags(
				metricsChan,
				statsTags,
			),
		),
	}

	opts = append(opts, getTLSOptsIfNeeded(tlsServerConf)...)

	oauthOpts, err := getOauthOptsIfNeeded(logger)
	if err != nil {
		logger.Fatal().AnErr("getOauthOptsIfNeeded", err).Msg("")
	}
	opts = append(opts, oauthOpts...)

	grpcServer := grpc.NewServer(opts...)
	methods.CreateAndRegisterServer(logger, grpcServer)

	logger.Debug().Str("service", "{{.ServiceName}}").
		Msg("starting grpc server")
	go grpcServer.Serve(lis)

	cancels, err = notifyZkOfRPCServerIfNeeded(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	zkCancels = append(
		zkCancels,
		cancels...,
	)

	if viper.GetBool("use_gateway_proxy") {
		logger.Debug().Str("service", "{{.ServiceName}}").
			Msg("starting gateway proxy")
		if err = startGatewayProxy(ctx, logger); err != nil {
			logger.Fatal().AnErr("startGatewayProxy", err).Msg("")
		}
	}

	cancels, err = notifyZkOfGatewayEndpointIfNeeded(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	zkCancels = append(
		zkCancels,
		cancels...,
	)

	s := <- sigChan
	logger.Info().Str("service", "{{.ServiceName}}") .
		Str("signal", s.String()).
		Msg("shutting down")
	cancelFunc()
	grpcServer.Stop()
}
