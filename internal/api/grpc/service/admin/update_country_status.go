package admin

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/country"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/guard"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/country"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s country.Service) UpdateCountryStatus(ctx context.Context, req *svc.UpdateCountryStatusRequest) (*svc.Country, error) {
	_, err := guard.AllowedRoles(ctx, req.Initiator, "create profile",
		roles.SuperUser, roles.Admin)
	if err != nil {
		return nil, err
	}

	countryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		return nil, problem.InvalidArgumentError(ctx, "invalid country id format", &errdetails.BadRequest_FieldViolation{
			Field:       "country_id",
			Description: "invalid UUID format for country ID",
		})
	}

	status, err := enum.ParseCountryStatus(req.Status)
	if err != nil {
		return nil, problem.InvalidArgumentError(ctx, "request id", &errdetails.BadRequest_FieldViolation{
			Field:       "status",
			Description: err.Error(),
		})
	}

	country, err := s.app.UpdateCountryStatus(ctx, countryID, status)
	if err != nil {
		return nil, err
	}

	return response.Country(country), nil
}
