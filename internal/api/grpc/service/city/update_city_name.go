package city

import (
	"context"

	svc "github.com/chains-lab/cities-proto/gen/go/svc/city"
	"github.com/chains-lab/cities-svc/internal/api/grpc/meta"
	"github.com/chains-lab/cities-svc/internal/api/grpc/responses"

	"github.com/chains-lab/cities-svc/internal/problems"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateCityName(ctx context.Context, req *svc.UpdateCityNameRequest) (*svc.City, error) {
	user := meta.User(ctx)

	var cityID uuid.UUID
	if user.Role == roles.Admin || user.Role == roles.SuperUser {
		var err error
		cityID, err = uuid.Parse(req.CityId)
		if err != nil {
			logger.Log(ctx).WithError(err).Error("invalid city ID format")

			return nil, problems.InvalidArgumentError(ctx, "city_id is invalid", &errdetails.BadRequest_FieldViolation{
				Field:       "city_id",
				Description: "invalid UUID format for city ID",
			})
		}
	} else {
		gov, err := s.OnlyCityAdmin(ctx, user.ID.String(), req.CityId, "update city name")
		if err != nil {
			return nil, err
		}

		cityID = gov.CityID
	}

	city, err := s.app.UpdateCityName(ctx, cityID, req.Name)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update city name")

		return nil, err
	}

	logger.Log(ctx).Infof("city name updated by user %s for city ID %s", user.ID, city.ID)

	return responses.City(city), nil
}
