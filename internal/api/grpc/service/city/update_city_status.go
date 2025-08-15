package city

import (
	"context"

	cityProto "github.com/chains-lab/cities-dir-proto/gen/go/svc/city"
	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/city"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/meta"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateCityStatus(ctx context.Context, req *svc.UpdateCityStatusRequest) (*cityProto.City, error) {
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
		gov, err := s.OnlyCityAdmin(ctx, user.ID.String(), req.CityId, "update city status")
		if err != nil {
			return nil, err
		}

		cityID = gov.CityID
	}

	cityStatus, err := enum.ParseCityStatus(req.Status)
	if err != nil {
		logger.Log(ctx).Error(err)

		return nil, problems.InvalidArgumentError(ctx, "city status is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "status",
			Description: err.Error(),
		})
	}

	city, err := s.app.UpdateCitiesStatus(ctx, cityID, cityStatus)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to update city status")

		return nil, err
	}

	logger.Log(ctx).Infof("city status updated by user %s for city ID %s", cityID, city.ID)

	return responses.City(city), nil
}
