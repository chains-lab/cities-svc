package grpc

import (
	"context"
	"fmt"
	"net"

	cityProto "github.com/chains-lab/cities-proto/gen/go/svc/city"
	countryProto "github.com/chains-lab/cities-proto/gen/go/svc/country"
	"github.com/chains-lab/cities-svc/internal/api/grpc/interceptors"
	"github.com/chains-lab/cities-svc/internal/api/grpc/service/city"
	"github.com/chains-lab/cities-svc/internal/api/grpc/service/country"
	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/config"

	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg config.Config, log logger.Logger, app *app.App) error {
	logInt := logger.UnaryLogInterceptor(log)
	requestId := interceptors.RequestID()
	userAuth := interceptors.UserJwtAuth(cfg.JWT.User.AccessToken.SecretKey)
	serviceAuth := interceptors.ServiceJwtAuth(cfg.JWT.Service.SecretKey)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logInt,
			requestId,
			serviceAuth,
			userAuth,
		),
	)

	cityProto.RegisterCityServiceServer(grpcServer, city.NewService(cfg, app))
	countryProto.RegisterCountryServiceServer(grpcServer, country.NewService(cfg, app))

	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	log.Infof("gRPC server listening on %s", lis.Addr())

	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- grpcServer.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down gRPC server …")
		grpcServer.GracefulStop()
		return nil
	case err := <-serveErrCh:
		return fmt.Errorf("gRPC Serve() exited: %w", err)
	}
}
