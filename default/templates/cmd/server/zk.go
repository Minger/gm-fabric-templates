package main

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"github.com/deciphernow/gm-fabric-go/gk"
)

type zkCancelFunc func()

func notifyZkOfMetricsIfNeeded(logger zerolog.Logger) ([]zkCancelFunc, error) {
	if !viper.GetBool("use_zk") {
		return nil, nil
	}

	host, err := checkAnnounceHost(viper.GetString("zk_announce_host"), logger)
	if err != nil {
		return nil, err
	}

	zk_servers := strings.Split(viper.GetString("zk_connection_string"), ",")

	logger.Info().Str("service", "{{.ServiceName}}").Msg("announcing metrics endpoint to zookeeper")
	cancel := gk.Announce(zk_servers, &gk.Registration{
		Path:   viper.GetString("zk_announce_path") + viper.GetString("metrics_dashboard_uri_path"),
		Host:   host,
		Status: gk.Alive,
		Port:   viper.GetInt("metrics_server_port"),
	})
	logger.Info().Str("service", "{{.ServiceName}}").Msg("Service successfully registered metrics endpoint to zookeeper")

	return []zkCancelFunc{cancel}, nil
}

func notifyZkOfRPCServerIfNeeded(logger zerolog.Logger) ([]zkCancelFunc, error) {
	if !viper.GetBool("use_zk") {
		return nil, nil
	}

	host, err := checkAnnounceHost(viper.GetString("zk_announce_host"), logger)
	if err != nil {
		return nil, err
	}

	zk_servers := strings.Split(viper.GetString("zk_connection_string"), ",")

	logger.Info().Str("service", "{{.ServiceName}}").Msg("announcing rpc endpoint to zookeeper")
	cancel := gk.Announce(zk_servers, &gk.Registration{
		Path:   viper.GetString("zk_announce_path") + "/rpc",
		Host:   host,
		Status: gk.Alive,
		Port:   viper.GetInt("grpc_server_port"),
	})
	logger.Info().Str("service", "{{.ServiceName}}").Msg("Service successfully registered rpc endpoint to zookeeper")

	return []zkCancelFunc{cancel}, nil
}

func notifyZkOfGatewayEndpointIfNeeded(logger zerolog.Logger) ([]zkCancelFunc, error) {
	if !(viper.GetBool("use_zk") && viper.GetBool("use_gateway_proxy")) {
		return nil, nil
	}

	host, err := checkAnnounceHost(viper.GetString("zk_announce_host"), logger)
	if err != nil {
		return nil, err
	}

	gatewayEndpoint := "http"
	if viper.GetBool("use_tls") {
		gatewayEndpoint = "https"
	}

	zk_servers := strings.Split(viper.GetString("zk_connection_string"), ",")

	logger.Info().Str("service", "{{.ServiceName}}").Msg("announcing gateway endpoint to zookeeper")

	cancel := gk.Announce(zk_servers, &gk.Registration{
		Path:   viper.GetString("zk_announce_path") + "/" + gatewayEndpoint,
		Host:   host,
		Status: gk.Alive,
		Port:   viper.GetInt("gateway_proxy_port"),
	})
	logger.Info().Str("service", "{{.ServiceName}}").Msg("announcing gateway endpoint to zookeeper")

	return []zkCancelFunc{cancel}, nil
}

func checkAnnounceHost(ah string, logger zerolog.Logger) (string, error) {
	if ah == "" {
		ip, err := gk.GetIP()
		if err != nil {
			logger.Error().AnErr("get_ip", err).Msg("Failed to get IP")
			return "", err
		}

		return ip, nil
	}
	return ah, nil
}
