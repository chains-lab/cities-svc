package app

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/cities-dir-svc/internal/app/models"
	"github.com/chains-lab/cities-dir-svc/internal/constant/enum"
	"github.com/chains-lab/cities-dir-svc/internal/dbx"
	"github.com/chains-lab/cities-dir-svc/internal/errx"
	"github.com/chains-lab/cities-dir-svc/internal/pagination"
	"github.com/google/uuid"
)

type formQ interface {
	New() dbx.FormToCreateCityQ

	Insert(ctx context.Context, input dbx.FormToCreateCityModel) error
	Get(ctx context.Context) (dbx.FormToCreateCityModel, error)
	Select(ctx context.Context) ([]dbx.FormToCreateCityModel, error)
	Update(ctx context.Context, input dbx.UpdateFormToCreateCityInput) error
	Delete(ctx context.Context) error

	FilterID(ID uuid.UUID) dbx.FormToCreateCityQ
	FilterStatus(status string) dbx.FormToCreateCityQ
	FilterInitiatorID(initiatorID uuid.UUID) dbx.FormToCreateCityQ
	FilterCountryID(countryID uuid.UUID) dbx.FormToCreateCityQ

	CityNameLike(name string) dbx.FormToCreateCityQ

	Page(limit, offset uint64) dbx.FormToCreateCityQ
	Count(ctx context.Context) (uint64, error)
}

type CreateFormInput struct {
	CityName     string
	CountryID    uuid.UUID
	InitiatorID  uuid.UUID
	ContactEmail string
	ContactPhone string
	Text         string
}

func (a App) CreateForm(ctx context.Context, input CreateFormInput) (models.Form, error) {
	ID := uuid.New()
	status := enum.FormToCreateCityStatusPending
	now := time.Now().UTC()

	form := dbx.FormToCreateCityModel{
		ID:           ID,
		Status:       status,
		CityName:     input.CityName,
		CountryID:    input.CountryID,
		InitiatorID:  input.InitiatorID,
		ContactEmail: input.ContactEmail,
		ContactPhone: input.ContactPhone,
		Text:         input.Text,
		UserRevID:    uuid.Nil, // No user reviewed yet
		CreatedAt:    now,
	}

	if err := a.formQ.New().Insert(ctx, form); err != nil {
		switch {
		default:
			return models.Form{}, errx.RaiseInternal(ctx, err)
		}
	}

	return formModel(form), nil
}

func (a App) AcceptForm(ctx context.Context, initiatorID, formID, adminID uuid.UUID) (models.Form, error) {
	form, err := a.formQ.New().FilterID(formID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Form{}, errx.RaiseFormNotFoundByID(ctx, err, formID)
		default:
			return models.Form{}, errx.RaiseInternal(ctx, err)
		}
	}

	trxErr := a.transaction(func(ctx context.Context) error {
		status := enum.FormToCreateCityStatusAccepted

		updateInput := dbx.UpdateFormToCreateCityInput{
			Status:    &status,
			UserRevID: &initiatorID,
		}

		if err := a.formQ.New().FilterID(formID).Update(ctx, updateInput); err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return errx.RaiseFormNotFoundByID(ctx, err, formID)
			default:
				return errx.RaiseInternal(ctx, err)
			}
		}

		city, err := a.CreateCity(ctx, CreateCityInput{
			CountryID: form.CountryID,
			Name:      form.CityName,
			Status:    enum.CityStatusSupported,
		})
		if err != nil {
			return err
		}

		_, err = a.CreateCityGov(ctx, city.ID, adminID, CreateCityGovInput{
			Role: enum.CityAdminRoleAdmin,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if trxErr != nil {
		return models.Form{}, trxErr
	}

	res, err := a.GetForm(ctx, formID)
	if err != nil {
		return models.Form{}, err
	}

	return res, nil
}

func (a App) RejectForm(ctx context.Context, formID uuid.UUID, reason string) (models.Form, error) {
	form, err := a.formQ.New().FilterID(formID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Form{}, errx.RaiseFormNotFoundByID(ctx, err, formID)
		default:
			return models.Form{}, errx.RaiseInternal(ctx, err)
		}
	}

	status := enum.FormToCreateCityStatusRejected
	updateInput := dbx.UpdateFormToCreateCityInput{
		Status: &status,
		Text:   &reason,
	}

	if err := a.formQ.New().FilterID(formID).Update(ctx, updateInput); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Form{}, errx.RaiseFormNotFoundByID(ctx, err, formID)
		default:
			return models.Form{}, errx.RaiseInternal(ctx, err)
		}
	}

	return formModel(form), nil
}

func (a App) GetForm(ctx context.Context, formID uuid.UUID) (models.Form, error) {
	form, err := a.formQ.New().FilterID(formID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Form{}, errx.RaiseFormNotFoundByID(ctx, err, formID)
		default:
			return models.Form{}, errx.RaiseInternal(ctx, err)
		}
	}

	return formModel(form), nil
}

type SearchFormsInput struct {
	Status      *string
	CountryID   *uuid.UUID
	CityName    *string
	InitiatorID *uuid.UUID
}

func (a App) SearchForms(ctx context.Context, input SearchFormsInput, pagPar pagination.Request, newFirst bool) ([]models.Form, pagination.Response, error) {
	limit, offset := pagination.CalculateLimitOffset(pagPar)

	query := a.formQ.New().Page(limit, offset)

	if input.Status != nil {
		query = query.FilterStatus(*input.Status)
	}

	if input.CountryID != nil {
		query = query.FilterCountryID(*input.CountryID)
	}

	if input.CityName != nil {
		query = query.CityNameLike(*input.CityName)
	}

	if input.InitiatorID != nil {
		query = query.FilterInitiatorID(*input.InitiatorID)
	}

	forms, err := query.Select(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []models.Form{}, pagination.Response{}, nil
		default:
			return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
		}
	}

	total, err := query.Count(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			total = 0
		default:
			return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
		}
	}

	res, pagResp := formsArray(forms, limit, offset, total)

	return res, pagResp, nil
}

func formModel(input dbx.FormToCreateCityModel) models.Form {
	return models.Form{
		ID:           input.ID,
		Status:       input.Status,
		CityName:     input.CityName,
		CountryID:    input.CountryID,
		InitiatorID:  input.InitiatorID,
		ContactEmail: input.ContactEmail,
		ContactPhone: input.ContactPhone,
		Text:         input.Text,
		CreatedAt:    input.CreatedAt,
	}
}

func formsArray(forms []dbx.FormToCreateCityModel, limit, offset, total uint64) ([]models.Form, pagination.Response) {
	res := make([]models.Form, len(forms))
	for i, form := range forms {
		res[i] = formModel(form)
	}

	pagResp := pagination.Response{
		Page:  offset/limit + 1,
		Size:  limit,
		Total: total,
	}

	return res, pagResp
}
