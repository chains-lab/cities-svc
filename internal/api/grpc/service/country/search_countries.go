package country

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/country"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
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

	countries, pag, err := s.app.SearchCountries(ctx, req.NameLike, status, pagination.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to search countries")

		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	return responses.CountriesList(countries, pag), nil
}
