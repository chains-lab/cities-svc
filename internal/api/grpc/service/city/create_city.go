package city

import (
	"context"

	svc "github.com/chains-lab/cities-proto/gen/go/svc/city"
	"github.com/chains-lab/cities-svc/internal/api/grpc/meta"
	"github.com/chains-lab/cities-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) CreateCity(ctx context.Context, req *svc.CreateCityRequest) (*svc.City, error) {
	user := meta.User(ctx)

	if user.Role != roles.Admin && user.Role != roles.SuperUser {
		logger.Log(ctx).Error("user does not have permission to create city")

		return nil, problems.UnauthenticatedError(ctx, "user does not have permission to create city")
	}

	countryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid country ID format")

		return nil, problems.InvalidArgumentError(ctx, "country_id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "country_id",
			Description: "invalid UUID format for country ID",
		})
	}

	city, err := s.app.CreateCity(ctx, app.CreateCityInput{
		Name:      req.Name,
		CountryID: countryID,
		Status:    req.Status,
	})
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to create city")

		return nil, problems.InternalError(ctx)
	}

	logger.Log(ctx).Infof("city created by %s, city ID: %s", user.ID, city.ID)

	return responses.City(city), nil
}
