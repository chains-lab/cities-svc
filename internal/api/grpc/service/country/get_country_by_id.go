package country

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/country"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
)

func (s Service) GetCountryById(ctx context.Context, req *svc.GetCountryByIdRequest) (*svc.Country, error) {
	countryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("invalid country ID format")

		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "id",
			Description: "invalid UUID format for country ID",
		})
	}

	country, err := s.app.GetCountryByID(ctx, countryID)
	if err != nil {
		logger.Log(ctx, RequestID(ctx)).WithError(err).Error("failed to get country by ID")

		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	logger.Log(ctx, RequestID(ctx)).Infof("retrieved country with ID %s", country.ID)

	return responses.Country(country), nil
}
