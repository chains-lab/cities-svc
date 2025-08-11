package country

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/country"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateCountryStatus(ctx context.Context, req *svc.UpdateCountryStatusRequest) (*svc.Country, error) {
	countryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		return nil, problems.InvalidArgumentError(ctx, "invalid country id format", &errdetails.BadRequest_FieldViolation{
			Field:       "country_id",
			Description: "invalid UUID format for country ID",
		})
	}

	status, err := enum.ParseCountryStatus(req.Status)
	if err != nil {
		return nil, problems.InvalidArgumentError(ctx, "request id", &errdetails.BadRequest_FieldViolation{
			Field:       "status",
			Description: err.Error(),
		})
	}

	country, err := s.app.UpdateCountryStatus(ctx, countryID, status)
	if err != nil {
		return nil, err
	}

	return responses.Country(country), nil
}
