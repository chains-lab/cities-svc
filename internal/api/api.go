package api

import (
	"context"
	"fmt"
	"net"

	cityProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/city"
	cityAdminProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/cityadmin"
	countryProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/country"
	countryAdminProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/countryadmin"
	formProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/form"
	formAdminProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/formadmin"
	cityGovProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/gov"
	cityGovAdminProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/govadmin"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/cityadmin"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/countryadmin"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/form"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/formadmin"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/govadmin"

	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/interceptor"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/city"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/country"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/gov"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/config"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg config.Config, log logger.Logger, app *app.App) error {
	logInterceptor := logger.UnaryLogInterceptor(log)
	authInterceptor := interceptor.Auth(cfg.JWT.Service.SecretKey)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logInterceptor, authInterceptor,
		),
	)

	cityProto.RegisterCityServiceServer(grpcServer, city.NewService(cfg, app))
	cityAdminProto.RegisterCityAdminServiceServer(grpcServer, cityadmin.NewService(cfg, app))
	cityGovProto.RegisterGovServiceServer(grpcServer, gov.NewService(cfg, app))
	cityGovAdminProto.RegisterGovAdminServiceServer(grpcServer, govadmin.NewService(cfg, app))
	countryProto.RegisterCountryServiceServer(grpcServer, country.NewService(cfg, app))
	countryAdminProto.RegisterCountryAdminServiceServer(grpcServer, countryadmin.NewService(cfg, app))
	formProto.RegisterFormServiceServer(grpcServer, form.NewService(cfg, app))
	formAdminProto.RegisterFormAdminServiceServer(grpcServer, formadmin.NewService(cfg, app))

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
