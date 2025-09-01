package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/constant"
	"github.com/chains-lab/cities-svc/internal/dbx"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type City struct {
	citiesQ dbx.CitiesQ

	slugRegexp *regexp.Regexp
	nameRegexp *regexp.Regexp
}

func NewCity(db *sql.DB) City {
	return City{
		citiesQ:    dbx.NewCitiesQ(db),
		slugRegexp: regexp.MustCompile(`^[a-z]+(-[a-z]+)*$`),
		nameRegexp: regexp.MustCompile(`^[A-Za-z -]+$`),
	}
}

func (c City) validateTimezone(tz string) error {
	if tz == "" {
		return errx.ErrorInvalidTimeZone.Raise(
			fmt.Errorf("timezone must not be empty"),
		)
	}
	_, err := time.LoadLocation(tz)
	if err != nil {
		return errx.ErrorInvalidTimeZone.Raise(
			fmt.Errorf("invalid timezone: %s", tz),
		)
	}
	return nil
}

func (c City) validatePoint(p orb.Point) error {
	lon, lat := p[0], p[1]

	if lon < -180 || lon > 180 {
		return errx.ErrorInvalidPoint.Raise(
			fmt.Errorf("invalid longitude: %.6f (must be between -180 and 180)", lon),
		)
	}
	if lat < -90 || lat > 90 {
		return errx.ErrorInvalidPoint.Raise(
			fmt.Errorf("invalid latitude: %.6f (must be between -90 and 90)", lat),
		)
	}
	return nil
}

func (c City) validateSlug(slug string) error {
	if slug == "" {
		return errx.ErrorInvalidSlug.Raise(
			fmt.Errorf("slug must not be empty"),
		)
	}
	if !c.slugRegexp.MatchString(slug) {
		return errx.ErrorInvalidSlug.Raise(
			fmt.Errorf("invalid slug: %s", slug),
		)
	}
	return nil
}

func (c City) validateName(name string) error {
	if name == "" {
		return errx.ErrorInvalidCityName.Raise(
			fmt.Errorf("city name must not be empty"),
		)
	}
	if !c.nameRegexp.MatchString(name) {
		return errx.ErrorInvalidCityName.Raise(
			fmt.Errorf("invalid city name: %s", name),
		)
	}
	return nil
}

type CreateCityParams struct {
	CountryID uuid.UUID
	Name      string
	Timezone  string
	Point     orb.Point
}

func (c City) Create(ctx context.Context, params CreateCityParams) (models.City, error) {
	err := c.validateTimezone(params.Timezone)
	if err != nil {
		return models.City{}, err
	}

	err = c.validatePoint(params.Point)
	if err != nil {
		return models.City{}, err
	}

	err = c.validateName(params.Name)
	if err != nil {
		return models.City{}, err
	}

	cityID := uuid.New()
	now := time.Now().UTC()

	resp := models.City{
		ID:        cityID,
		CountryID: params.CountryID,
		Status:    constant.CityStatusCommunity,
		Timezone:  params.Timezone,
		CreatedAt: now,
		UpdatedAt: now,
	}

	stmt := dbx.City{
		ID:        cityID,
		CountryID: params.CountryID,
		Point:     params.Point,
		Status:    constant.CityStatusCommunity,
		Name:      params.Name,
		Timezone:  params.Timezone,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = c.citiesQ.Insert(ctx, stmt)
	if err != nil {
		return models.City{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to creating city: %w", err),
		)
	}

	return resp, nil
}

// Read methods for city

func (c City) GetByID(ctx context.Context, cityID uuid.UUID) (models.City, error) {
	city, err := c.citiesQ.New().FilterID(cityID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, errx.ErrorCityNotFound.Raise(
				fmt.Errorf("сity not found by id: %s, cause: %w", cityID, err),
			)
		default:
			return models.City{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city by id: %s, cause: %w", cityID, err),
			)
		}
	}

	return cityFromDb(city), nil
}

func (c City) GetByRadius(ctx context.Context, point orb.Point, radius uint64) (models.City, error) {
	city, err := c.citiesQ.New().
		FilterWithinRadiusMeters(point, radius).
		OrderByNearest(point, true).
		Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, errx.ErrorCityNotFound.Raise(
				fmt.Errorf("nearest city not found, cause: %w", err),
			)
		default:
			return models.City{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get nearest city, cause: %w", err),
			)
		}
	}

	return cityFromDb(city), nil
}

func (c City) GetBySlug(ctx context.Context, slug string) (models.City, error) {
	city, err := c.citiesQ.New().FilterSlug(slug).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.City{}, errx.ErrorCityNotFound.Raise(
				fmt.Errorf("city not found by slug: %s, cause: %w", slug, err),
			)
		default:
			return models.City{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get city by slug, cause: %w", err),
			)
		}
	}

	return cityFromDb(city), nil
}

type SelectCityParams struct {
	Name      *string
	Status    []string
	CountryID *uuid.UUID
	Point     *orb.Point
}

func (c City) SelectCities(
	ctx context.Context,
	params SelectCityParams,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.City, pagi.Response, error) {
	if pag.Page == 0 {
		pag.Page = 1
	}
	if pag.Size == 0 {
		pag.Size = 20
	}
	if pag.Size > 100 {
		pag.Size = 100
	}

	limit := pag.Size + 1 // +1 чтобы определить наличие next
	offset := (pag.Page - 1) * pag.Size

	query := c.citiesQ.New()

	if params.Name != nil {
		query = query.FilterNameLike(*params.Name)
	}

	for _, s := range params.Status {
		err := constant.CheckCityStatus(s)
		if err != nil {
			return nil, pagi.Response{}, errx.ErrorInvalidCityStatus.Raise(
				fmt.Errorf("failed to invalid city status: %s, cause: %w", s, err),
			)
		}
	}
	if len(params.Status) > 0 {
		query = query.FilterStatus(params.Status...)
	}
	if params.CountryID != nil {
		query = query.FilterCountryID(*params.CountryID)
	}

	for _, s := range sort {
		switch s.Field {
		case "name":
			query = query.OrderByAlphabetical(s.Ascend)
		case "distance":
			if params.Point != nil {
				query = query.OrderByNearest(*params.Point, s.Ascend)
			}
		}
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to count cities, cause: %w", err),
		)
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to search cities, cause: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	cities := make([]models.City, 0, len(rows))
	for _, city := range rows {
		cities = append(cities, cityFromDb(city))
	}

	return cities, pagi.Response{
		Size:  pag.Size,
		Page:  pag.Page,
		Total: total,
	}, nil
}

type UpdateCityParams struct {
	CountryID *uuid.UUID
	Point     *orb.Point
	Status    *string
	Name      *string
	Icon      *string
	Slug      *string
	Timezone  *string

	UpdatedAt time.Time
}

func (c City) UpdateOne(ctx context.Context, cityID uuid.UUID, params UpdateCityParams) error {
	_, err := c.GetByID(ctx, cityID)
	if err != nil {
		return err
	}

	if params.CountryID == nil && params.Point == nil && params.Status == nil &&
		params.Name == nil && params.Icon == nil && params.Slug == nil && params.Timezone == nil {
		return nil
	}

	stmt := dbx.UpdateCityParams{}
	if params.CountryID != nil {
		stmt.CountryID = params.CountryID
	}
	if params.Point != nil {
		err = c.validatePoint(*params.Point)
		if err != nil {
			return err
		}
		stmt.Point = params.Point
	}
	if params.Status != nil {
		err = constant.CheckCityStatus(*params.Status)
		if err != nil {
			return errx.ErrorInvalidCityStatus.Raise(
				fmt.Errorf("failed to invalid city status, cause: %s", err),
			)
		}
		stmt.Status = params.Status
	}
	if params.Name != nil {
		err = c.validateName(*params.Name)
		if err != nil {
			return err
		}
		stmt.Name = params.Name
	}
	if params.Icon != nil && *params.Icon != "" {
		stmt.Icon = &sql.NullString{String: *params.Icon, Valid: true}
	} else if params.Icon != nil && *params.Icon == "" {
		stmt.Icon = &sql.NullString{String: "", Valid: false}
	}
	if params.Slug != nil && *params.Slug != "" {
		err = c.validateSlug(*params.Slug)
		if err != nil {
			return err
		}

		stmt.Slug = &sql.NullString{String: *params.Slug, Valid: true}
	} else if params.Slug != nil && *params.Slug == "" {
		stmt.Slug = &sql.NullString{String: "", Valid: false}
	}
	if params.Timezone != nil {
		err = c.validateTimezone(*params.Timezone)
		if err != nil {
			return err
		}
		stmt.Timezone = params.Timezone
	}

	stmt.UpdatedAt = params.UpdatedAt

	err = c.citiesQ.New().FilterID(cityID).Update(ctx, stmt)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update city, cause: %w", err),
		)
	}

	return nil
}

type UpdateCitiesFilters struct {
	CountryID *uuid.UUID
	Status    []string
}

type UpdateCitiesParams struct {
	Status   *string
	Timezone *string

	UpdatedAt time.Time
}

func (c City) UpdateMany(ctx context.Context, filters UpdateCitiesFilters, params UpdateCitiesParams) error {
	query := c.citiesQ.New()
	if filters.CountryID != nil {
		query = query.FilterCountryID(*filters.CountryID)
	}
	if filters.Status != nil {
		for _, s := range filters.Status {
			err := constant.CheckCityStatus(s)
			if err != nil {
				return errx.ErrorInvalidCityStatus.Raise(
					fmt.Errorf("failed to invalid city status: %s, cause: %w", s, err),
				)
			}
		}
		query = query.FilterStatus(filters.Status...)
	}

	if params.Status == nil && params.Timezone == nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("0 filters provided for update cities"),
		)
	}

	stmt := dbx.UpdateCityParams{
		UpdatedAt: params.UpdatedAt,
	}

	if params.Status != nil {
		err := constant.CheckCityStatus(*params.Status)
		if err != nil {
			return errx.ErrorInvalidCityStatus.Raise(
				fmt.Errorf("failed to invalid city status: %s, cause: %w", *params.Status, err),
			)
		}
		stmt.Status = params.Status
	}
	if params.Timezone != nil {
		err := c.validateTimezone(*params.Timezone)
		if err != nil {
			return err
		}
		stmt.Timezone = params.Timezone
	}

	err := query.Update(ctx, stmt)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update cities by country_id, cause: %w", err),
		)
	}

	return nil
}

func cityFromDb(c dbx.City) models.City {
	res := models.City{
		ID:        c.ID,
		CountryID: c.CountryID,
		Status:    c.Status,
		Name:      c.Name,
		Timezone:  c.Timezone,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	if c.Icon.Valid {
		res.Icon = &c.Icon.String
	}
	if c.Slug.Valid {
		res.Slug = &c.Slug.String
	}

	return res
}
