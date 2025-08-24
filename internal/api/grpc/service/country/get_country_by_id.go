package country

import (
	"context"

	svc "github.com/chains-lab/cities-proto/gen/go/svc/country"
	"github.com/chains-lab/cities-svc/internal/api/grpc/responses"

	"github.com/chains-lab/cities-svc/internal/problems"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetCountryById(ctx context.Context, req *svc.GetCountryByIdRequest) (*svc.Country, error) {
	countryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid country ID format")

		return nil, problems.InvalidArgumentError(ctx, "country id is invalid format", &errdetails.BadRequest_FieldViolation{
			Field:       "id",
			Description: "invalid UUID format for country ID",
		})
	}

	country, err := s.app.GetCountryByID(ctx, countryID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get country by ID")

		return nil, err
	}

	return responses.Country(country), nil
}
