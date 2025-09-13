package apptest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func CreateAndActivateCountry(s Setup, t *testing.T, name string) models.Country {
	ctx := context.Background()

	ukr, err := s.app.CreateCountry(ctx, name)
	if err != nil {
		t.Fatalf("CreateCountry: %v", err)
	}

	ukr, err = s.app.SetCountryStatusSupported(ctx, ukr.ID)
	if err != nil {
		t.Fatalf("SetCountryStatusSupported: %v", err)
	}
	if ukr.Status != enum.CountryStatusSupported {
		t.Errorf("expected country status 'supported', got '%s'", ukr.Status)
	}

	return ukr
}

func TestCreateCities(t *testing.T) {
	s, err := newSetup(t)
	cleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	//usa := CreateAndActivateCountry(s, t, "USA")
	//ukr, err = s.app.SetCountryStatusDeprecated(ctx, ukr.ID)
	//if err != nil {
	//	t.Fatalf("SetCountryStatusDeprecated: %v", err)
	//}

	kyiv, err := s.app.CreateCity(ctx, app.CreateCityParams{
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
	cleanDb(t)

	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	usa := CreateAndActivateCountry(s, t, "USA")
	usa, err = s.app.SetCountryStatusDeprecated(ctx, usa.ID)
	if err != nil {
		t.Fatalf("SetCountryStatusDeprecated: %v", err)
	}

	_, err = s.app.CreateCity(ctx, app.CreateCityParams{
		Name:      "New York",
		CountryID: usa.ID,
		Point:     [2]float64{30.5234, 50.4501},
		Timezone:  "America/New_York",
	})
	if !errors.Is(err, errx.ErrorCountryNotSupported) {
		t.Fatalf("expected error %v, got %v", errx.ErrorCountryNotSupported, err)
	}
}

func TestCreateCityInNonExistentCountry(t *testing.T) {
	s, err := newSetup(t)
	cleanDb(t)

	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	_, err = s.app.CreateCity(ctx, app.CreateCityParams{
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
	cleanDb(t)

	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")

	_, err = s.app.CreateCity(ctx, app.CreateCityParams{
		Name:      "Kyiv",
		CountryID: ukr.ID,
		Point:     [2]float64{30.5234, 50.4501},
		Timezone:  "Invalid/Timezone",
	})
	if !errors.Is(err, errx.ErrorInvalidTimeZone) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInvalidTimeZone, err)
	}
}

func TestCreateCityWithInvalidPoint(t *testing.T) {
	s, err := newSetup(t)
	cleanDb(t)

	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")

	_, err = s.app.CreateCity(ctx, app.CreateCityParams{
		Name:      "Kyiv",
		CountryID: ukr.ID,
		Point:     [2]float64{200.0, 50.4501}, // Invalid longitude
		Timezone:  "Europe/Kyiv",
	})
	if !errors.Is(err, errx.ErrorInvalidPoint) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInvalidPoint, err)
	}

	_, err = s.app.CreateCity(ctx, app.CreateCityParams{
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
	cleanDb(t)

	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}
	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")

	_, err = s.app.CreateCity(ctx, app.CreateCityParams{
		Name:      " ", // Empty name
		CountryID: ukr.ID,
		Point:     [2]float64{30.5234, 50.4501},
		Timezone:  "Europe/Kyiv",
	})
	if !errors.Is(err, errx.ErrorInvalidCityName) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInvalidCityName, err)
	}
}
