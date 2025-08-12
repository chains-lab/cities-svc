package citygov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) RefuseCityGovRight(ctx context.Context, req *svc.RefuseCityGovRightRequest) (*emptypb.Empty, error) {
	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid city ID format")

		return nil, problem.InvalidArgumentError(ctx, "invalid city id format", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	initiator, err := s.OnlyGov(ctx, req.Initiator, cityID, "refuse city government rights")
	if err != nil {
		return nil, err
	}

	if initiator.Role == enum.CityAdminRoleAdmin {
		logger.Log(ctx).Error("city admin try to refuse own admin rights")

		return nil, problem.PermissionDeniedError(ctx, "city admin cannot transfer own admin rights, but u can transfer to another user")
	}

	err = s.app.RefuseOwnCityGovRights(ctx, cityID, initiator.ID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to transfer city ownership")

		return nil, err
	}

	return &emptypb.Empty{}, nil
}
