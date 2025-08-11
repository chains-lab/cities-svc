package citygov

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/citygov"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problems"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/responses"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) ListUserCitiesAdmins(ctx context.Context, req *svc.ListUserCitiesAdminsRequest) (*svc.ListCitiesAdmins, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, problems.InvalidArgumentError(ctx, "user id is invalid format", &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "invalid UUID format for user ID",
		})
	}

	citiesAdmins, pag, err := s.app.GetUserCitiesAdmins(ctx, userID, pagination.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		return nil, err
	}

	return responses.CitiesAdminsList(citiesAdmins, pag), nil
}
