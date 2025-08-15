package gov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/gov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/meta"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) CreateCityGovAdmin(ctx context.Context, req *svc.CreateCityGovAdminRequest) (*svc.CityGov, error) {
	user := meta.User(ctx)

	if user.Role != roles.Admin && user.Role != roles.SuperUser {
		logger.Log(ctx).Warnf("user %s with role %s tried to create city admin, but only Admin or SuperUser can do this",
			user.ID, user.Role)

		return nil, problems.PermissionDeniedError(ctx, "only Admin or SuperUser can create city admin")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid user ID format")

		return nil, problems.InvalidArgumentError(ctx, "user id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid city ID format")

		return nil, problems.InvalidArgumentError(ctx, "city id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	cityAdmin, err := s.app.CreateCityGovAdmin(ctx, cityID, userID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to create city admin")

		return nil, err
	}

	logger.Log(ctx).Infof("city admin created by user %s for city ID %s", user.ID, cityAdmin.CityID)

	return responses.CityAdmin(cityAdmin), nil
}
