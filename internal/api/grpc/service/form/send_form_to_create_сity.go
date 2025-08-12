package form

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/svc/form"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/guard"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/logger"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) SendFormToCreateCity(ctx context.Context, req *svc.SendFormToCreateCityRequest) (*svc.FormToCreateCity, error) {
	initiatorID, err := guard.AllowedRoles(ctx, req.Initiator, "send form to create city", roles.User)
	if err != nil {
		return &svc.FormToCreateCity{}, err
	}

	countryId, err := uuid.Parse(req.CountryId)
	if err != nil {
		logger.Log(ctx).WithError(err).Error("invalid country ID format")

		return &svc.FormToCreateCity{}, problem.InvalidArgumentError(ctx, "country id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "country_id",
			Description: "invalid UUID format for country ID",
		})
	}

	form, err := s.app.CreateForm(ctx, app.CreateFormInput{
		CityName:     req.CityName,
		CountryID:    countryId,
		InitiatorID:  initiatorID,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
		Text:         req.Text,
	})
	if err != nil {
		logger.Log(ctx).WithError(err).Error("failed to send form to create city")
		return &svc.FormToCreateCity{}, err
	}

	logger.Log(ctx).Infof("form to create city sent by user %s for country %s", initiatorID, req.CountryId)

	return response.Form(form), nil
}
