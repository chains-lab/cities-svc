package formadmin

import (
	"context"

	"github.com/chains-lab/cities-dir-proto/gen/go/svc/form"
	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/formadmin"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/guard"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) AcceptFormToCreateCity(ctx context.Context, req *svc.AcceptFormToCreateCityRequest) (*form.FormToCreateCity, error) {
	initiatorID, err := guard.AllowedRoles(ctx, req.Initiator, "accept form to create city",
		roles.SuperUser, roles.Admin)
	if err != nil {
		return nil, err
	}

	AdminForCityID, err := uuid.Parse(req.AdminId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid admin ID format")

		return nil, problem.InvalidArgumentError(ctx, "admin_id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "admin_id",
			Description: "invalid UUID format for admin ID",
		})
	}

	FormID, err := uuid.Parse(req.FormId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid form ID format")

		return nil, problem.InvalidArgumentError(ctx, "form_id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "form_id",
			Description: "invalid UUID format for form ID",
		})
	}

	formToCreateCity, err := s.app.AcceptForm(ctx, initiatorID, FormID, AdminForCityID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("add form to create city failed")

		return nil, problem.InvalidArgumentError(ctx, "add form to create city failed", &errdetails.BadRequest_FieldViolation{
			Field:       "form",
			Description: err.Error(),
		})
	}

	logger.Log(ctx).Infof("form to create city accepted by %s, form ID: %s", initiatorID, FormID)

	return response.Form(formToCreateCity), nil
}
