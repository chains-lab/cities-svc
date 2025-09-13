package apptest

import (
	"context"
	"testing"

	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func TestGetCities(t *testing.T) {
	s, err := newSetup(t)
	cleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	cities, err := s.app.GetCityByID(ctx, kyiv.ID)
	if err != nil {
		t.Fatalf("GetCityByID: %v", err)
	}
	if cities.ID != kyiv.ID {
		t.Fatalf("GetCityByID: expected city ID %v, got %v", kyiv.ID, cities.ID)
	}
}

func TestGetCitiesBySlug(t *testing.T) {
	s, err := newSetup(t)
	cleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	cities, err := s.app.UpdateCity(ctx, kyiv.ID, uuid.New(), roles.SuperUser, app.UpdateCityParams{
		Slug: func(s string) *string { return &s }("kyiv"),
	})
	if err != nil {
		t.Fatalf("UpdateCity: %v", err)
	}

	cityBySlug, err := s.app.GetCityBySlug(ctx, "kyiv")
	if err != nil {
		t.Fatalf("GetCityBySlug: %v", err)
	}
	if cityBySlug.ID != cities.ID {
		t.Fatalf("GetCityBySlug: expected city ID %v, got %v", cities.ID, cityBySlug.ID)
	}
}

func TestListCities(t *testing.T) {
	s, err := newSetup(t)
	cleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	_ = CreateCity(s, t, ukr.ID, "Kyiv")
	_ = CreateCity(s, t, ukr.ID, "Lviv")

	usa := CreateAndActivateCountry(s, t, "USA")
	_ = CreateCity(s, t, usa.ID, "New York")

	cities, _, err := s.app.ListCities(ctx, app.FilterListCitiesParams{}, pagi.Request{
		Page: 1,
		Size: 10,
	}, []pagi.SortField{})
	if err != nil {
		t.Fatalf("ListCities: %v", err)
	}
	for _, c := range cities {
		t.Logf("City: %s, CountryID: %s", c.Name, c.CountryID)
	}
}
