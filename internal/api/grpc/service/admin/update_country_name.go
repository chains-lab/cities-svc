package admin

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/country"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/guard"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/service/country"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s country.Service) UpdateCountryName(ctx context.Context, req *svc.UpdateCountryNameRequest) (*svc.Country, error) {
	_, err := guard.AllowedRoles(ctx, req.Initiator, "create profile",
		roles.SuperUser, roles.Admin)
	if err != nil {
		return nil, err
	}

	countryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid country ID format")

		return nil, problem.InvalidArgumentError(ctx, "invalid country id format", &errdetails.BadRequest_FieldViolation{
			Field:       "id",
			Description: "invalid UUID format for country ID",
		})
	}

	country, err := s.app.UpdateCountryName(ctx, countryID, req.Name)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to create country")

		return nil, err
	}

	return response.Country(country), nil
}
