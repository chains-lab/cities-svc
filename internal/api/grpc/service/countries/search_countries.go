package countries

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/countries"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/renderer"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
)

func (s Service) SearchCountries(ctx context.Context, req *svc.SearchCountriesRequest) (*svc.CountriesList, error) {
	status, ok := enum.ParseCountryStatus(req.Status)
	if !ok {
		logger.Log(ctx, RequestID(ctx)).Errorf("invalid country status provided: %s", req.Status)

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "status",
			Description: "invalid country status provided",
		})
	}

	countries, err := s.methods.SearchCountries(ctx, req.NameLike, status, req.Pagination.Limit, req.Pagination.Page)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to search countries")

		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	return renderer.CountriesList(countries), nil
}
