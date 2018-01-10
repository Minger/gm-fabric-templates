package main

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"google.golang.org/grpc"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"github.com/deciphernow/gm-fabric-go/oauth"
)

func putOauthInCtxIfNeeded(ctx context.Context) context.Context {
	if viper.GetBool("use_oauth") {
		return oauth.ContextWithOptions(
			ctx,
			oauth.WithProvider(viper.GetString("oauth_provider")),
			oauth.WithClientID(viper.GetString("oauth_client_id")),
		)
	}
	return ctx
}

func getOauthOptsIfNeeded(logger zerolog.Logger) ([]grpc.ServerOption, error) {
	var err error

	if !viper.GetBool("use_oauth") {
		return nil, nil
	}

	provider := viper.GetString("oauth_provider")
	clientID := viper.GetString("oauth_client_id")

	logger.Debug().Str("service", "{{.ServiceName}}").
		Str("oauth_provider", provider).
		Str("oauth_client_id", clientID).
		Msg("loading OAuth config")

	interceptor, err := oauth.NewOauthInterceptor(
		oauth.WithProvider(provider),
		oauth.WithClientID(clientID),
	)
	if err != nil {
		return nil, errors.Wrap(err, "oauth.NewOauthInterceptor")
	}

	return []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(interceptor)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(interceptor)),
	}, nil
}

