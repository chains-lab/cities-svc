package gov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/gov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) RefuseCityGovRight(ctx context.Context, req *svc.RefuseCityGovRightRequest) (*emptypb.Empty, error) {
	initiator, err := s.OnlyGov(ctx, req.Initiator.UserId, req.CityId, "refuse city government rights")
	if err != nil {
		return nil, err
	}

	if initiator.Role == enum.CityAdminRoleAdmin {
		logger.Log(ctx).Error("city admin try to refuse own admin rights")

		return nil, problem.PermissionDeniedError(ctx, "city admin cannot transfer own admin rights, but u can transfer to another user")
	}

	err = s.app.RefuseOwnCityGovRights(ctx, initiator.CityID, initiator.ID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to transfer city ownership")

		return nil, err
	}

	return &emptypb.Empty{}, nil
}
