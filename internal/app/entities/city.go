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
	"github.com/jackc/pgx/v5"
	"github.com/paulmach/orb"
)

type City struct {
	citiesQ dbx.CitiesQ

	slugRegexp *regexp.Regexp
	nameRegexp *regexp.Regexp
}

func CreateCityEntity(db *sql.DB) City {
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
	Status    string
	Name      string
	Icon      *string
	Slug      *string
	Timezone  string
	Point     orb.Point
}

func (c City) Create(ctx context.Context, params CreateCityParams) (models.City, error) {
	err := constant.CheckCityStatus(params.Status)
	if err != nil {
		return models.City{}, errx.ErrorInvalidCityStatus.Raise(
			fmt.Errorf("failed to invalid city status, cause: %w", err),
		)
	}

	err = c.validateTimezone(params.Timezone)
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
		Status:    params.Status,
		Timezone:  params.Timezone,
		CreatedAt: now,
		UpdatedAt: now,
	}

	stmt := dbx.City{
		ID:        cityID,
		CountryID: params.CountryID,
		Point:     params.Point,
		Status:    params.Status,
		Name:      params.Name,
		Timezone:  params.Timezone,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if params.Slug != nil {
		_, err = c.citiesQ.New().FilterSlug(*params.Slug).Get(ctx)
		if err == nil {
			return models.City{}, errx.ErrorCityAlreadyExists.Raise(
				fmt.Errorf("failed to city already exists with slug: %v", params.Slug),
			)
		}

		stmt.Slug = sql.NullString{String: *params.Slug, Valid: true}
		resp.Slug = params.Slug
	}

	if params.Icon != nil {
		stmt.Icon = sql.NullString{String: *params.Icon, Valid: true}
		resp.Icon = params.Icon
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
		case errors.Is(err, pgx.ErrNoRows):
			return models.City{}, errx.ErrorCityNotFound.Raise(
				fmt.Errorf("ity not found by id: %s, cause: %w", cityID, err),
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

type CitySearchParams struct {
	Name    *string
	Status  []string
	Country *uuid.UUID
	Point   *orb.Point
}

func (c City) SearchCities(
	ctx context.Context,
	params CitySearchParams,
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
	if params.Country != nil {
		query = query.FilterCountryID(*params.Country)
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

type UpdateCityStatusParams struct {
	CountryID *uuid.UUID
	Point     *orb.Point
	Status    *string
	Name      *string
	Icon      *string
	Slug      *string
	Timezone  *string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c City) Update(ctx context.Context, cityID uuid.UUID, params UpdateCityStatusParams) (models.City, error) {
	city, err := c.GetByID(ctx, cityID)
	if err != nil {
		return models.City{}, err
	}

	stmt := dbx.UpdateCityParams{}
	if params.CountryID != nil {
		stmt.CountryID = params.CountryID
	}
	if params.Point != nil {
		err = c.validatePoint(*params.Point)
		if err != nil {
			return models.City{}, err
		}
		stmt.Point = params.Point
	}
	if params.Status != nil {
		err = constant.CheckCityStatus(*params.Status)
		if err != nil {
			return models.City{}, errx.ErrorInvalidCityStatus.Raise(
				fmt.Errorf("failed to invalid city status, cause: %s", err),
			)
		}
		stmt.Status = params.Status
	}
	if params.Name != nil {
		err = c.validateName(*params.Name)
		if err != nil {
			return models.City{}, err
		}
		stmt.Name = params.Name
	}
	if params.Icon != nil {
		stmt.Icon = &sql.NullString{String: *params.Icon, Valid: true}
	} else if *city.Icon == "" {
		stmt.Icon = &sql.NullString{String: "", Valid: false}
	}
	if params.Slug != nil {
		err = c.validateSlug(*params.Slug)
		if err != nil {
			return models.City{}, err
		}

		if city.Slug == nil || *city.Slug != *params.Slug {
			_, err = c.citiesQ.New().FilterSlug(*params.Slug).Get(ctx)
			if err == nil {
				return models.City{}, errx.ErrorCityAlreadyExists.Raise(
					fmt.Errorf("failed to city already exists with slug: %v", params.Slug),
				)
			}
		}

		stmt.Slug = &sql.NullString{String: *params.Slug, Valid: true}
	} else if *city.Slug == "" {
		// if slug is empty string, we set it to NULL in DB
		stmt.Slug = &sql.NullString{String: "", Valid: false}
	}
	if params.Timezone != nil {
		err = c.validateTimezone(*params.Timezone)
		if err != nil {
			return models.City{}, err
		}
		stmt.Timezone = params.Timezone
	}

	stmt.UpdatedAt = time.Now().UTC()

	return c.GetByID(ctx, cityID)
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
