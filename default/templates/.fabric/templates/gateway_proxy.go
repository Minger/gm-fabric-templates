package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/deciphernow/gm-fabric-go/httputil"
	"github.com/deciphernow/gm-fabric-go/metrics/httpmeta"
	"github.com/deciphernow/gm-fabric-go/metrics/proxymeta"
	"github.com/deciphernow/gm-fabric-go/middleware"
	"github.com/deciphernow/gm-fabric-go/rpcutil"
	"github.com/deciphernow/gm-fabric-go/tlsutil"

	pb "{{"{{.PBImport}}"}}"	
)

func startGatewayProxy(ctx context.Context, logger zerolog.Logger) error {
	var listener net.Listener
	var err error

	mux := runtime.NewServeMux(proxymeta.MetaOption(), runtime.WithIncomingHeaderMatcher(rpcutil.MatchHTTPHeaders))
	var handler http.Handler = mux 

	if viper.GetBool("verbose_logging") {
		stack := middleware.Chain(
			middleware.MiddlewareFunc(hlog.NewHandler(logger)),
			middleware.MiddlewareFunc(hlog.AccessHandler(func(r *http.Request, status int, size int, duration time.Duration) {
				hlog.FromRequest(r).Info().
					Str("method", r.Method).
					Str("path", r.URL.String()).
					Int("status", status).
					Int("size", size).
					Dur("duration", duration).
					Msg("Access")
			})),
			middleware.MiddlewareFunc(hlog.UserAgentHandler("user_agent")),
		)
		handler = stack.Wrap(handler)
	}

	if err = registerClient(ctx, logger, mux); err != nil {
		return errors.Wrap(err, "registerClient")
	}

	proxyAddress := fmt.Sprintf(
		"%s:%d",
		viper.GetString("gateway_proxy_host"),
		viper.GetInt("gateway_proxy_port"),
	)
	if viper.GetBool("gateway_use_tls") {
		var tlsServerConf *tls.Config

		tlsServerConf, err = tlsutil.BuildServerTLSConfig(
			viper.GetString("ca_cert_path"),
			viper.GetString("server_cert_path"),
			viper.GetString("server_key_path"),
		)
		if err != nil {
			return errors.Wrap(err, "tlsutil.BuildServerTLSConfig")
		}
		listener, err = tls.Listen("tcp", proxyAddress, tlsServerConf)
		if err != nil {
			return errors.Wrap(err, "tls.Listen failed")
		}
	} else {
		listener, err = net.Listen("tcp", proxyAddress)
		if err != nil {
			return errors.Wrap(err, "tls.Listen failed")
		}
	}

	logger.Debug().Str("service", "{{"{{.ServiceName}}"}}").
		Str("host", viper.GetString("gateway_proxy_host")).
		Int("port", viper.GetInt("gateway_proxy_port")).
		Msg("starting gateway proxy server")

	// compose the meta handler
	handler = httpmeta.Handler(handler)

	// compose the grpc array handler
	if viper.GetBool("gatweway_serve_anonymous_arrays") {
		logger.Debug()Str("service", "{{"{{.ServiceName}}"}}").
			Msg("serving anonymous arrays")
		handler = httputil.Handler(logger, handler)
	}

	go http.Serve(listener, handler)

	return nil
}

func registerClient(
	ctx context.Context, 
	logger zerolog.Logger, 
	mux *runtime.ServeMux,
) error {
	var err error 

	var clientOpts []grpc.DialOption
	if viper.GetBool("grpc_use_tls") {
		var creds credentials.TransportCredentials
		var tlsClientConf *tls.Config

		tlsClientConf, err = tlsutil.NewTLSClientConfig(
			viper.GetString("ca_cert_path"),
			viper.GetString("server_cert_path"),
			viper.GetString("server_key_path"),
			viper.GetString("server_cert_name"),
		)
		if err != nil {
			return errors.Wrap(err, "tlsutil.NewTLSClientConfig")
		}

		creds = credentials.NewTLS(tlsClientConf)
		clientOpts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	} else {
		clientOpts = []grpc.DialOption{grpc.WithInsecure()}
	}

	err = pb.Register{{"{{.ServerInterfaceName}}"}}HandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf(
			"%s:%d",
			viper.GetString("grpc_server_host"),
			viper.GetInt("grpc_server_port"),
		),
		clientOpts,
	)
	if err != nil {
		return errors.Wrap(err, "pb.Register{{"{{.ServerInterfaceName}}"}}HandlerFromEndpoint")
	}

	return nil
}

