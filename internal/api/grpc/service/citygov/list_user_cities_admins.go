package citygov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
	"github.com/google/uuid"
)

func (s Service) ListUserCitiesAdmins(ctx context.Context, req *svc.ListUserCitiesAdminsRequest) (*svc.ListCitiesAdmins, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, responses.InvalidArgumentError(ctx, RequestID(ctx), responses.Violation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	citiesAdmins, pag, err := s.app.GetUserCitiesAdmins(ctx, userID, pagination.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		return nil, responses.AppError(ctx, RequestID(ctx), err)
	}

	return responses.CitiesAdminsList(citiesAdmins, pag), nil
}
