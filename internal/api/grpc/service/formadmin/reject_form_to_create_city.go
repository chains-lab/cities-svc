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

func (s Service) RejectFormToCreateCity(ctx context.Context, req *svc.RejectFormToCreateCityRequest) (*form.FormToCreateCity, error) {
	_, err := guard.AllowedRoles(ctx, req.Initiator, "decline form to create city",
		roles.SuperUser, roles.Admin)
	if err != nil {
		return nil, err
	}

	FormID, err := uuid.Parse(req.FormId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid form ID format")

		return nil, problem.InvalidArgumentError(ctx, "form_id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "form_id",
			Description: "invalid UUID format for form ID",
		})
	}

	formToCreateCity, err := s.app.RejectForm(ctx, FormID, req.Reason)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("reject form to create city failed")

		return nil, problem.InvalidArgumentError(ctx, "reject form to create city failed", &errdetails.BadRequest_FieldViolation{
			Field:       "form",
			Description: err.Error(),
		})
	}

	logger.Log(ctx).Infof("form to create city rejected, form ID: %s", FormID)

	return response.Form(formToCreateCity), nil
}
