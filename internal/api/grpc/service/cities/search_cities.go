package cities

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/cities"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/renderer"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) SearchCities(ctx context.Context, req *svc.SearchCitiesRequest) (*svc.CitiesList, error) {
	CountryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid country ID format")

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "country_id",
			Description: "invalid UUID format for country ID",
		})
	}

	cities, err := s.methods.SearchCityInCountry(ctx, req.NameLike, CountryID, req.Pagination.Page, req.Pagination.Limit)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to search cities")
		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	return renderer.CitiesList(cities), nil
}
