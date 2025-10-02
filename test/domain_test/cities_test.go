package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/city"
	"github.com/chains-lab/cities-svc/test"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

func CreateAndActivateCountry(s Setup, t *testing.T, name string) models.Country {
	ctx := context.Background()

	ukr, err := s.domain.country.Create(ctx, name)
	if err != nil {
		t.Fatalf("CreateCity: %v", err)
	}

	ukr, err = s.domain.country.UpdateStatus(ctx, ukr.ID, enum.CountryStatusSupported)
	if err != nil {
		t.Fatalf("SetCountryStatusSupported: %v", err)
	}
	if ukr.Status != enum.CountryStatusSupported {
		t.Errorf("expected country status 'supported', got '%s'", ukr.Status)
	}

	return ukr
}

func CreateCity(s Setup, t *testing.T, countryID uuid.UUID, name string) models.City {
	ctx := context.Background()

	kyiv, err := s.domain.city.Create(ctx, city.CreateParams{
		Name:      name,
		CountryID: countryID,
		Point:     [2]float64{30.5234, 50.4501}, // Longitude, Latitude
		Timezone:  "Europe/Kyiv",
	})
	if err != nil {
		t.Fatalf("CreateCity: %v", err)
	}

	return kyiv
}

func TestCreateCities(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")

	kyiv, err := s.domain.city.Create(ctx, city.CreateParams{
		Name:      "Kyiv",
		CountryID: ukr.ID,
		Point:     [2]float64{30.5234, 50.4501},
		Timezone:  "Europe/Kyiv",
	})
	if err != nil {
		t.Fatalf("CreateCity: %v", err)
	}

	if kyiv.Name != "Kyiv" {
		t.Errorf("expected city name 'Kyiv', got '%s'", kyiv.Name)
	}
}

func TestCreateCityInUnsupportedCountry(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)

	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	usa := CreateAndActivateCountry(s, t, "USA")
	usa, err = s.domain.country.UpdateStatus(ctx, usa.ID, enum.CountryStatusDeprecated)
	if err != nil {
		t.Fatalf("SetCountryStatusDeprecated: %v", err)
	}

	_, err = s.domain.city.Create(ctx, city.CreateParams{
		Name:      "New York",
		CountryID: usa.ID,
		Point:     [2]float64{30.5234, 50.4501},
		Timezone:  "America/New_York",
	})
	if !errors.Is(err, errx.ErrorCountryIsNotSupported) {
		t.Fatalf("expected error %v, got %v", errx.ErrorCountryIsNotSupported, err)
	}
}

func TestCreateCityInNonExistentCountry(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)

	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	_, err = s.domain.city.Create(ctx, city.CreateParams{
		Name:      "New York",
		CountryID: uuid.New(),
		Point:     [2]float64{30.5234, 50.4501},
		Timezone:  "America/New_York",
	})
	if !errors.Is(err, errx.ErrorCountryNotFound) {
		t.Fatalf("expected error %v, got %v", errx.ErrorCountryNotFound, err)
	}
}

func TestCreateCityWithInvalidTimezone(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)

	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")

	_, err = s.domain.city.Create(ctx, city.CreateParams{
		Name:      "New York",
		CountryID: ukr.ID,
		Point:     [2]float64{30.5234, 50.4501},
		Timezone:  "America/Pidor", // Invalid timezone
	})
	if !errors.Is(err, errx.ErrorInvalidTimeZone) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInvalidTimeZone, err)
	}
}

func TestCreateCityWithInvalidPoint(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)

	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")

	_, err = s.domain.city.Create(ctx, city.CreateParams{
		Name:      "Kyiv",
		CountryID: ukr.ID,
		Point:     [2]float64{200.0, 50.4501}, // Invalid longitude
		Timezone:  "Europe/Kyiv",
	})
	if !errors.Is(err, errx.ErrorInvalidPoint) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInvalidPoint, err)
	}

	_, err = s.domain.city.Create(ctx, city.CreateParams{
		Name:      "Kyiv",
		CountryID: ukr.ID,
		Point:     [2]float64{30.5234, 100.0}, // Invalid latitude
		Timezone:  "Europe/Kyiv",
	})
	if !errors.Is(err, errx.ErrorInvalidPoint) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInvalidPoint, err)
	}
}

func TestCreateCityWithInvalidName(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)

	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}
	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")

	_, err = s.domain.city.Create(ctx, city.CreateParams{
		Name:      " ", // Empty name
		CountryID: ukr.ID,
		Point:     [2]float64{30.5234, 50.4501},
		Timezone:  "Europe/Kyiv",
	})
	if !errors.Is(err, errx.ErrorInvalidCityName) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInvalidCityName, err)
	}
}

func TestGetCities(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	cities, err := s.domain.city.GetByID(ctx, kyiv.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if cities.ID != kyiv.ID {
		t.Fatalf("GetByID: expected city ID %v, got %v", kyiv.ID, cities.ID)
	}
}

func TestGetCitiesBySlug(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	cities, err := s.domain.city.Update(ctx, kyiv.ID, city.UpdateParams{
		Slug: func(s string) *string { return &s }("kyiv"),
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	cityBySlug, err := s.domain.city.GetBySlug(ctx, "kyiv")
	if err != nil {
		t.Fatalf("GetBySlug: %v", err)
	}
	if cityBySlug.ID != cities.ID {
		t.Fatalf("GetBySlug: expected city ID %v, got %v", cities.ID, cityBySlug.ID)
	}
}

func TestListCities(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	_ = CreateCity(s, t, ukr.ID, "Kyiv")
	_ = CreateCity(s, t, ukr.ID, "Lviv")

	usa := CreateAndActivateCountry(s, t, "USA")
	_ = CreateCity(s, t, usa.ID, "New York")

	cities, err := s.domain.city.Filter(ctx, city.FilterParams{}, 1, 10)
	if err != nil {
		t.Fatalf("ListCities: %v", err)
	}
	for _, c := range cities.Data {
		t.Logf("City: %s, CountryID: %s", c.Name, c.CountryID)
	}
}

func TestUpdateCities(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	cityUpdZone := "Europe/Kyiv"
	cityUpdName := "Kyiv Updated"
	point := orb.Point{30.5234, 50.4501}
	kyiv, err = s.domain.city.Update(ctx, kyiv.ID, city.UpdateParams{
		Name:     &cityUpdName,
		Timezone: &cityUpdZone,
		Point:    &point,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
}

func TestUpdateCityFailsForNonExistentCity(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	cityUpdZone := "Europe/Kyiv"
	cityUpdName := "Kyiv Updated"
	point := orb.Point{30.5234, 50.4501}
	_, err = s.domain.city.Update(ctx, uuid.New(), city.UpdateParams{
		Name:     &cityUpdName,
		Timezone: &cityUpdZone,
		Point:    &point,
	})
	if !errors.Is(err, errx.ErrorCityNotFound) {
		t.Fatalf("Update: %v", err)
	}
}

func TestUpdateCityFailsForInvalidTimezone(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	cityUpdZone := "Europe/Invalid"
	cityUpdName := "Kyiv Updated"
	point := orb.Point{30.5234, 50.4501}
	kyiv, err = s.domain.city.Update(ctx, kyiv.ID, city.UpdateParams{
		Name:     &cityUpdName,
		Timezone: &cityUpdZone, // Invalid timezone
		Point:    &point,
	})
	if !errors.Is(err, errx.ErrorInvalidTimeZone) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInvalidTimeZone, err)
	}
}

func TestUpdateCityFailsForInvalidPoint(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	cityUpdZone := "Europe/Kyiv"
	cityUpdName := "Kyiv Updated"
	point := orb.Point{301.5234, 501.4501}
	_, err = s.domain.city.Update(ctx, kyiv.ID, city.UpdateParams{
		Name:     &cityUpdName,
		Timezone: &cityUpdZone,
		Point:    &point,
	})
	if !errors.Is(err, errx.ErrorInvalidPoint) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInvalidPoint, err)
	}
}

func TestUpdateCityFailsForInvalidName(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	cityUpdZone := "Europe/Kyiv"
	cityUpdName := "Kyiv 12121 Updated"
	point := orb.Point{01.5234, 01.4501}
	_, err = s.domain.city.Update(ctx, kyiv.ID, city.UpdateParams{
		Name:     &cityUpdName,
		Timezone: &cityUpdZone,
		Point:    &point,
	})
	if !errors.Is(err, errx.ErrorInvalidCityName) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInvalidCityName, err)
	}
}

func TestSetCityStatus(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	kyiv, err = s.domain.city.UpdateStatus(ctx, kyiv.ID, enum.CityStatusOfficial)
	if err != nil {
		t.Fatalf("SetCityStatusOfficial: %v", err)
	}
	if kyiv.Status != enum.CityStatusOfficial {
		t.Errorf("expected city status 'official', got '%s'", kyiv.Status)
	}

	kyiv, err = s.domain.city.UpdateStatus(ctx, kyiv.ID, enum.CityStatusDeprecated)
	if err != nil {
		t.Fatalf("SetCityStatusDeprecated: %v", err)
	}
	if kyiv.Status != enum.CityStatusDeprecated {
		t.Errorf("expected city status 'deprecated', got '%s'", kyiv.Status)
	}

	_, err = s.domain.moder.CreateInvite(context.Background(), enum.CityAdminRoleHead, kyiv.ID, time.Hour)
	if !errors.Is(err, errx.ErrorCityIsNotSupported) {
		t.Fatalf("expected error when creating mayor invite for deprecated city, got: %v", err)
	}
}

func TestSetCityStatusInUnsupportedCountry(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")
	lviv := CreateCity(s, t, ukr.ID, "Lviv")

	ukr, err = s.domain.country.UpdateStatus(ctx, ukr.ID, enum.CountryStatusDeprecated)
	if err != nil {
		t.Fatalf("SetCountryStatusDeprecated: %v", err)
	}
	if ukr.Status != enum.CountryStatusDeprecated {
		t.Errorf("expected country status 'deprecated', got '%s'", ukr.Status)
	}

	kyiv, err = s.domain.city.GetByID(ctx, kyiv.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if kyiv.Status != enum.CityStatusDeprecated {
		t.Errorf("expected city status 'deprecated' after country deprecated, got '%s'", kyiv.Status)
	}
	lviv, err = s.domain.city.GetByID(ctx, lviv.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if lviv.Status != enum.CityStatusDeprecated {
		t.Errorf("expected city status 'deprecated' after country deprecated, got '%s'", lviv.Status)
	}

	_, err = s.domain.city.UpdateStatus(ctx, kyiv.ID, enum.CityStatusOfficial)
	if !errors.Is(err, errx.ErrorCountryIsNotSupported) {
		t.Fatalf("expected error when setting city status in deprecated country, got: %v", err)
	}

	ukr, err = s.domain.country.UpdateStatus(ctx, ukr.ID, enum.CountryStatusSupported)
	if err != nil {
		t.Fatalf("Update country status to active: %v", err)
	}
	if ukr.Status != enum.CountryStatusSupported {
		t.Errorf("expected country status 'supported', got '%s'", ukr.Status)
	}

	kyiv, err = s.domain.city.GetByID(ctx, kyiv.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if kyiv.Status != enum.CityStatusDeprecated {
		t.Errorf("expected city status 'deprecated' after country supported, got '%s'", kyiv.Status)
	}
	lviv, err = s.domain.city.GetByID(ctx, lviv.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if lviv.Status != enum.CityStatusDeprecated {
		t.Errorf("expected city status 'deprecated' after country supported, got '%s'", lviv.Status)
	}

	kyiv, err = s.domain.city.UpdateStatus(ctx, kyiv.ID, enum.CityStatusOfficial)
	if err != nil {
		t.Fatalf("SetCityStatusOfficial: %v", err)
	}
	if kyiv.Status != enum.CityStatusOfficial {
		t.Errorf("expected city status 'official', got '%s'", kyiv.Status)
	}
}
