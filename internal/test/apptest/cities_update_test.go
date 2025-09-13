package apptest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

func CreateCity(s Setup, t *testing.T, countryID uuid.UUID, name string) models.City {
	ctx := context.Background()

	kyiv, err := s.app.CreateCity(ctx, app.CreateCityParams{
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

func TestUpdateCities(t *testing.T) {
	s, err := newSetup(t)
	cleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	kyiv, err = s.app.UpdateCity(ctx, kyiv.ID, uuid.New(), roles.SuperUser, app.UpdateCityParams{
		Name:     func(s string) *string { return &s }("Kyiv Updated"),
		Timezone: func(s string) *string { return &s }("Europe/Kyiv"),
		Point:    func(p orb.Point) *orb.Point { return &p }(orb.Point{30.5234, 50.4501}),
	})
	if err != nil {
		t.Fatalf("UpdateCity: %v", err)
	}
}

func TestUpdateCityFailsForNonExistentCity(t *testing.T) {
	s, err := newSetup(t)
	cleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	_, err = s.app.UpdateCity(ctx, uuid.New(), uuid.New(), roles.SuperUser, app.UpdateCityParams{
		Name:     func(s string) *string { return &s }("Kyiv Updated"),
		Timezone: func(s string) *string { return &s }("Europe/Kyiv"),
		Point:    func(p orb.Point) *orb.Point { return &p }(orb.Point{30.5234, 50.4501}),
	})
	if !errors.Is(err, errx.ErrorCityNotFound) {
		t.Fatalf("UpdateCity: %v", err)
	}
}

func TestUpdateCityFailsForInvalidTimezone(t *testing.T) {
	s, err := newSetup(t)
	cleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	_, err = s.app.UpdateCity(ctx, kyiv.ID, uuid.New(), roles.SuperUser, app.UpdateCityParams{
		Name:     func(s string) *string { return &s }("Kyiv Updated"),
		Timezone: func(s string) *string { return &s }("Invalid/Timezone"),
		Point:    func(p orb.Point) *orb.Point { return &p }(orb.Point{30.5234, 50.4501}),
	})
	if !errors.Is(err, errx.ErrorInvalidTimeZone) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInvalidTimeZone, err)
	}
}

func TestUpdateCityFailsForInvalidPoint(t *testing.T) {
	s, err := newSetup(t)
	cleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	_, err = s.app.UpdateCity(ctx, kyiv.ID, uuid.New(), roles.SuperUser, app.UpdateCityParams{
		Name:     func(s string) *string { return &s }("Kyiv Updated"),
		Timezone: func(s string) *string { return &s }("Europe/Kyiv"),
		Point:    func(p orb.Point) *orb.Point { return &p }(orb.Point{200.0, 100.0}), // Invalid point
	})
	if !errors.Is(err, errx.ErrorInvalidPoint) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInvalidPoint, err)
	}
}

func TestUpdateCityFailsForInvalidName(t *testing.T) {
	s, err := newSetup(t)
	cleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	_, err = s.app.UpdateCity(ctx, kyiv.ID, uuid.New(), roles.SuperUser, app.UpdateCityParams{
		Name:     func(s string) *string { return &s }(" "),
		Timezone: func(s string) *string { return &s }("Europe/Kyiv"),
		Point:    func(p orb.Point) *orb.Point { return &p }(orb.Point{30.5234, 50.4501}),
	})
	if !errors.Is(err, errx.ErrorInvalidCityName) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInvalidCityName, err)
	}
}
