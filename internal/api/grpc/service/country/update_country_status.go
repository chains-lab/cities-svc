package country

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/country"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/google/uuid"
)

func (s Service) UpdateCountryStatus(ctx context.Context, req *svc.UpdateCountryStatusRequest) (*svc.Country, error) {
	countryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "country_id",
			Description: "invalid UUID format for country ID",
		})
	}

	status, err := enum.ParseCountryStatus(req.Status)
	if err != nil {
		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "status",
			Description: err.Error(),
		})
	}

	country, err := s.methods.UpdateCountryStatus(ctx, countryID, status)
	if err != nil {
		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	return responses.Country(country), nil
}
