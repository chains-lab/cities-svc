package countries

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/countries"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/renderer"
	"github.com/chains-lab/cities-dir-svc/internal/enum"
	"github.com/google/uuid"
)

func (s Service) UpdateCountryStatus(ctx context.Context, req *svc.UpdateCountryStatusRequest) (*svc.Country, error) {
	countryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "country_id",
			Description: "invalid UUID format for country ID",
		})
	}

	status, ok := enum.ParseCountryStatus(req.Status)
	if !ok {
		return nil, renderer.InvalidArgumentError(ctx, RequestID(ctx), renderer.Violation{
			Field:       "status",
			Description: "invalid country status provided",
		})
	}

	country, err := s.methods.UpdateCountryStatus(ctx, countryID, status)
	if err != nil {
		return nil, renderer.AppError(ctx, RequestID(ctx), err)
	}

	return renderer.Country(country), nil
}
