package city

import (
	"context"

	svc "github.com/chains-lab/cities-dir-proto/gen/go/city"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/problem"
	"github.com/chains-lab/cities-dir-svc/internal/api/grpc/response"
	"github.com/chains-lab/cities-dir-svc/internal/app"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) SearchFormToCreateCity(ctx context.Context, req *svc.SearchFormToCreateCityRequest) (*svc.FormToCreateCityList, error) {
	var input app.SearchFormsInput

	if req.CityName != nil {
		input.CityName = req.CityName
	}

	if req.CountryId != nil {
		countryId, err := uuid.Parse(*req.CountryId)
		if err != nil {
			return nil, problem.InvalidArgumentError(ctx, "country id is invalid", &errdetails.BadRequest_FieldViolation{
				Field:       "country_id",
				Description: "invalid UUID format for country ID",
			})
		}
		input.CountryID = &countryId
	}

	if req.Status != nil {
		status, err := enum.ParseCityStatus(*req.Status)
		if err != nil {
			return nil, problem.InvalidArgumentError(ctx, "status is invalid", &errdetails.BadRequest_FieldViolation{
				Field:       "status",
				Description: err.Error(),
			})
		}
		input.Status = &status
	}

	if req.InitiatorId != nil {
		initiatorId, err := uuid.Parse(*req.InitiatorId)
		if err != nil {
			return nil, problem.InvalidArgumentError(ctx, "initiator id is invalid", &errdetails.BadRequest_FieldViolation{
				Field:       "initiator_id",
				Description: "invalid UUID format for initiator ID",
			})
		}
		input.InitiatorID = &initiatorId
	}

	forms, pagResp, err := s.app.SearchForms(ctx, input, pagination.Request{
		Page: req.Pag.Page,
		Size: req.Pag.Size,
	}, true)
	if err != nil {
		return nil, err
	}

	return response.FormList(forms, pagResp), nil
}
