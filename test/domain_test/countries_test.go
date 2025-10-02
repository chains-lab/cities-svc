package domain_test

import (
	"context"
	"testing"

	"github.com/chains-lab/cities-svc/internal/domain/services/country"
	"github.com/chains-lab/cities-svc/test"
	"github.com/chains-lab/enum"
)

func TestCountries(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	test.CleanDb(t)

	ctx := context.Background()

	ukr, err := s.domain.country.Create(ctx, "Ukraine")
	if err != nil {
		t.Fatalf("CreateCity: %v", err)
	}

	if ukr.Name != "Ukraine" {
		t.Errorf("expected country name 'Ukraine', got '%s'", ukr.Name)
	}

	var UkrNewName = "Ukraine Updated"
	ukr, err = s.domain.country.Update(ctx, ukr.ID, country.UpdateParams{
		Name: &UkrNewName,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	ukr2, err := s.domain.country.GetByID(ctx, ukr.ID)
	if err != nil {
		t.Fatalf("GetCountry: %v", err)
	}
	if ukr2.Name != UkrNewName {
		t.Errorf("expected country name '%s', got '%s'", UkrNewName, ukr2.Name)
	}

	usa, err := s.domain.country.Create(ctx, "USA")
	if err != nil {
		t.Fatalf("CreateCity: %v", err)
	}

	if usa.Name != "USA" {
		t.Errorf("expected country name 'USA', got '%s'", usa.Name)
	}

	countries, err := s.domain.country.Filter(ctx, country.FilterParams{}, 1, 10)
	if err != nil {
		t.Fatalf("ListCountries: %v", err)
	}

	if len(countries.Data) != 2 {
		t.Errorf("expected 2 country, got %d", len(countries.Data))
	}
	if countries.Total != 2 {
		t.Errorf("expected total 2 country, got %d", countries.Total)
	}

	if ukr.Status != enum.CountryStatusUnsupported {
		t.Errorf("expected country status 'unsuported', got '%s'", ukr.Status)
	}
	ukr, err = s.domain.country.UpdateStatus(ctx, ukr.ID, enum.CountryStatusDeprecated)
	if err != nil {
		t.Fatalf("SetCountryStatusDeprecated: %v", err)
	}
	if ukr.Status != enum.CountryStatusDeprecated {
		t.Errorf("expected country status 'deprecated', got '%s'", ukr.Status)
	}

	ukr, err = s.domain.country.UpdateStatus(ctx, ukr.ID, enum.CountryStatusSupported)
	if err != nil {
		t.Fatalf("SetCountryStatusSupported: %v", err)
	}
	if ukr.Status != enum.CountryStatusSupported {
		t.Errorf("expected country status 'supported', got '%s'", ukr.Status)
	}

	usa, err = s.domain.country.UpdateStatus(ctx, usa.ID, enum.CountryStatusDeprecated)
	if err != nil {
		t.Fatalf("SetCountryStatusDeprecated: %v", err)
	}
	if usa.Status != enum.CountryStatusDeprecated {
		t.Errorf("expected country status 'deprecated', got '%s'", usa.Status)
	}

	usa, err = s.domain.country.UpdateStatus(ctx, usa.ID, enum.CountryStatusSupported)
	if err != nil {
		t.Fatalf("SetCountryStatusSupported: %v", err)
	}
	if usa.Status != enum.CountryStatusSupported {
		t.Errorf("expected country status 'supported', got '%s'", usa.Status)
	}
}
