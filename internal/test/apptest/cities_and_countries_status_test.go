package apptest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

func TestSetCityStatus(t *testing.T) {
	s, err := newSetup(t)
	cleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	kyiv, err = s.app.SetCityStatusOfficial(ctx, kyiv.ID)
	if err != nil {
		t.Fatalf("SetCityStatusOfficial: %v", err)
	}
	if kyiv.Status != enum.CityStatusOfficial {
		t.Errorf("expected city status 'official', got '%s'", kyiv.Status)
	}

	userID := uuid.New()

	mayor := CreateMayor(s, t, kyiv.ID, userID)
	mayor, err = s.app.GetCityMayor(ctx, kyiv.ID)
	if err != nil {
		t.Fatalf("GetCityMayor: %v", err)
	}
	if mayor.UserID != userID {
		t.Errorf("expected mayor user ID '%s', got '%s'", userID, mayor.UserID)
	}

	kyiv, err = s.app.SetCityStatusDeprecated(ctx, kyiv.ID)
	if err != nil {
		t.Fatalf("SetCityStatusDeprecated: %v", err)
	}
	if kyiv.Status != enum.CityStatusDeprecated {
		t.Errorf("expected city status 'deprecated', got '%s'", kyiv.Status)
	}

	_, err = s.app.CreateInviteMayor(context.Background(), kyiv.ID, uuid.New(), roles.SuperUser)
	if !errors.Is(err, errx.ErrorCannotCreateMayorInviteForNotOfficialCity) {
		t.Fatalf("expected error when creating mayor invite for deprecated city, got: %v", err)
	}
}

func TestSetCityStatusInUnsupportedCountry(t *testing.T) {
	s, err := newSetup(t)
	cleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")
	lviv := CreateCity(s, t, ukr.ID, "Lviv")

	ukr, err = s.app.SetCountryStatusDeprecated(ctx, ukr.ID)
	if err != nil {
		t.Fatalf("SetCountryStatusDeprecated: %v", err)
	}
	if ukr.Status != enum.CountryStatusDeprecated {
		t.Errorf("expected country status 'deprecated', got '%s'", ukr.Status)
	}

	kyiv, err = s.app.GetCityByID(ctx, kyiv.ID)
	if err != nil {
		t.Fatalf("GetCityByID: %v", err)
	}
	if kyiv.Status != enum.CityStatusDeprecated {
		t.Errorf("expected city status 'deprecated' after country deprecated, got '%s'", kyiv.Status)
	}
	lviv, err = s.app.GetCityByID(ctx, lviv.ID)
	if err != nil {
		t.Fatalf("GetCityByID: %v", err)
	}
	if lviv.Status != enum.CityStatusDeprecated {
		t.Errorf("expected city status 'deprecated' after country deprecated, got '%s'", lviv.Status)
	}

	_, err = s.app.SetCityStatusOfficial(ctx, kyiv.ID)
	if !errors.Is(err, errx.ErrorCannotUpdateCityStatusInUnsupportedCountry) {
		t.Fatalf("expected error when setting city status in deprecated country, got: %v", err)
	}

	ukr, err = s.app.SetCountryStatusSupported(ctx, ukr.ID)
	if err != nil {
		t.Fatalf("SetCountryStatusActive: %v", err)
	}
	if ukr.Status != enum.CountryStatusSupported {
		t.Errorf("expected country status 'supported', got '%s'", ukr.Status)
	}

	kyiv, err = s.app.GetCityByID(ctx, kyiv.ID)
	if err != nil {
		t.Fatalf("GetCityByID: %v", err)
	}
	if kyiv.Status != enum.CityStatusDeprecated {
		t.Errorf("expected city status 'deprecated' after country supported, got '%s'", kyiv.Status)
	}
	lviv, err = s.app.GetCityByID(ctx, lviv.ID)
	if err != nil {
		t.Fatalf("GetCityByID: %v", err)
	}
	if lviv.Status != enum.CityStatusDeprecated {
		t.Errorf("expected city status 'deprecated' after country supported, got '%s'", lviv.Status)
	}

	kyiv, err = s.app.SetCityStatusOfficial(ctx, kyiv.ID)
	if err != nil {
		t.Fatalf("SetCityStatusOfficial: %v", err)
	}
	if kyiv.Status != enum.CityStatusOfficial {
		t.Errorf("expected city status 'official', got '%s'", kyiv.Status)
	}
}
