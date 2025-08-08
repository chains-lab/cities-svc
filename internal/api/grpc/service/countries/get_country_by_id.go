package countries

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/countries"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/renderer"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) GetCountryById(ctx context.Context, req *svc.GetCountryByIdRequest) (*svc.Country, error) {
	countryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid country ID format")

		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "id",
			Description: "invalid UUID format for country ID",
		})
	}

	country, err := s.methods.GetCountryByID(ctx, countryID)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to get country by ID")

		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	logger.Log(ctx, RequestID(ctx)).Infof("retrieved country with ID %s", country.ID)

	return renderer.Country(country), nil
}
