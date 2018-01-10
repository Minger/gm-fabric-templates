package main

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"github.com/deciphernow/gm-fabric-go/tlsutil"
)

func buildMetricsTLSConfigIfNeeded(logger zerolog.Logger) (*tls.Config, error) {
	if !viper.GetBool("metrics_use_tls") {
		logger.Debug().Str("service", "{{.ServiceName}}").Msg("not using metrics TLS")
		return nil, nil
	}

	tlsConf, err := createConfig(logger)
	if err != nil {
		return nil, err
	}

	return tlsConf, nil
}

func buildServerTLSConfigIfNeeded(logger zerolog.Logger) (*tls.Config, error) {
	if !viper.GetBool("grpc_use_tls") {
		logger.Debug().Str("service", "{{.ServiceName}}").Msg("not using grpc server TLS")
		return nil, nil
	}

	tlsConf, err := createConfig(logger)
	if err != nil {
		return nil, err
	}

	return tlsConf, nil
}

func createConfig(logger zerolog.Logger) (*tls.Config, error) {
	logger.Debug().Str("service", "{{.ServiceName}}").
		Str("ca_cert_path", viper.GetString("ca_cert_path")).
		Str("server_cert_path", viper.GetString("server_cert_path")).
		Str("server_key_path", viper.GetString("server_key_path")).
		Msg("loading TLS config")

	tlsConf, err := tlsutil.BuildServerTLSConfig(
		viper.GetString("ca_cert_path"),
		viper.GetString("server_cert_path"),
		viper.GetString("server_key_path"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "tlsutil.BuildServerTLSConfig")
	}

	return tlsConf, nil
}

func getTLSOptsIfNeeded(tlsServerConf *tls.Config) []grpc.ServerOption {
	if viper.GetBool("grpc_use_tls") {
		return []grpc.ServerOption{grpc.Creds(credentials.NewTLS(tlsServerConf))}
	}

	return nil
}
