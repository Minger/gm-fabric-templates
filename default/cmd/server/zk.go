package main

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"github.com/deciphernow/gm-fabric-go/gk"
)

type zkCancelFunc func()

func notifyZkOfMetricsIfNeeded(logger zerolog.Logger) []zkCancelFunc {
	if !viper.GetBool("use_zk") {
		return nil
	}

	logger.Info().Str("service", "{{.ServiceName}}").Msg("announcing metrics endpoint to zookeeper")
	cancel := gk.Announce(viper.GetStringSlice("zk_connection_string"), &gk.Registration{
		Path:   viper.GetString("zk_announce_path") + viper.GetString("metrics_uri_path"),
		Host:   viper.GetString("zk_announce_host"),
		Status: gk.Alive,
		Port:   viper.GetInt("metrics_server_port"),
	})
	logger.Info().Str("service", "{{.ServiceName}}").Msg("Service successfully registered metrics endpoint to zookeeper")

	return []zkCancelFunc{cancel}
}

func notifyZkOfRPCServerIfNeeded(logger zerolog.Logger) []zkCancelFunc {
	if !viper.GetBool("use_zk") {
		return nil
	}

	logger.Info().Str("service", "{{.ServiceName}}").Msg("announcing rpc endpoint to zookeeper")
	cancel := gk.Announce(viper.GetStringSlice("zk_connection_string"), &gk.Registration{
		Path:   viper.GetString("zk_announce_path") + "/rpc",
		Host:   viper.GetString("zk_announce_host"),
		Status: gk.Alive,
		Port:   viper.GetInt("grpc_server_port"),
	})
	logger.Info().Str("service", "{{.ServiceName}}").Msg("Service successfully registered rpc endpoint to zookeeper")

	return []zkCancelFunc{cancel}
}

func notifyZkOfGatewayEndpointIfNeeded(logger zerolog.Logger) []zkCancelFunc {
	if !(viper.GetBool("use_zk") && viper.GetBool("use_gateway_proxy")) {
		return nil
	}

	gatewayEndpoint := "http"
	if viper.GetBool("use_tls") {
		gatewayEndpoint = "https"
	}

	logger.Info().Str("service", "{{.ServiceName}}").Msg("announcing gateway endpoint to zookeeper")

	cancel := gk.Announce(viper.GetStringSlice("zk_connection_string"), &gk.Registration{
		Path:   viper.GetString("zk_announce_path") + "/" + gatewayEndpoint,
		Host:   viper.GetString("zk_announce_host"),
		Status: gk.Alive,
		Port:   viper.GetInt("gateway_proxy_port"),
	})
	logger.Info().Str("service", "{{.ServiceName}}").Msg("announcing gateway endpoint to zookeeper")

	return []zkCancelFunc{cancel}
}
