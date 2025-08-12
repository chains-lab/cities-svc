package countryadmin

import (
	"context"

	countryProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/country"
	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/countryadmin"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/guard"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateCountryName(ctx context.Context, req *svc.UpdateCountryNameRequest) (*countryProto.Country, error) {
	initiatorID, err := guard.AllowedRoles(ctx, req.Initiator, "create profile",
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

	logger.Log(ctx).Infof("country name updated by user %s for country ID %s", initiatorID, country.ID)

	return response.Country(country), nil
}
