package form

import (
	"context"
	"fmt"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/form"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetFormToCreateCity(ctx context.Context, req *svc.GetFormToCreateCityRequest) (*svc.FormToCreateCity, error) {
	formID, err := uuid.Parse(req.FormId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid form ID format")

		return nil, problem.InvalidArgumentError(ctx, fmt.Sprint("form_id is invalid"), &errdetails.BadRequest_FieldViolation{
			Field:       "form_id",
			Description: "invalid UUID format for form ID",
		})
	}

	form, err := s.app.GetForm(ctx, formID)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to get form to create city")

		return nil, problem.InternalError(ctx)
	}

	return response.Form(form), nil
}
