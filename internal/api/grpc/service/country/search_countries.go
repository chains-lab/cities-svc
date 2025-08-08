package country

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/country"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
)

func (s Service) SearchCountries(ctx context.Context, req *svc.SearchCountriesRequest) (*svc.CountriesList, error) {
	status, err := enum.ParseCountryStatus(req.Status)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).Error(err)

		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "status",
			Description: err.Error(),
		})
	}

	countries, err := s.app.SearchCountries(ctx, req.NameLike, status, req.Pagination.Limit, req.Pagination.Page)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to search countries")

		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	return responses.CountriesList(countries), nil
}
