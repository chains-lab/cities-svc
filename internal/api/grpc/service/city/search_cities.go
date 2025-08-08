package city

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/city"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) SearchCities(ctx context.Context, req *svc.SearchCitiesRequest) (*svc.CitiesList, error) {
	CountryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid country ID format")

		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "country_id",
			Description: "invalid UUID format for country ID",
		})
	}

	cities, err := s.app.SearchCityInCountry(ctx, req.NameLike, CountryID, req.Pagination.Page, req.Pagination.Limit)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to search cities")
		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	return responses.CitiesList(cities), nil
}
