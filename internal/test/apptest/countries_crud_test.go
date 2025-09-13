package apptest

import (
	"context"
	"testing"

	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
)

func TestCRUDCountries(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	ukr, err := s.app.CreateCountry(ctx, "Ukraine")
	if err != nil {
		t.Fatalf("CreateCountry: %v", err)
	}

	if ukr.Name != "Ukraine" {
		t.Errorf("expected country name 'Ukraine', got '%s'", ukr.Name)
	}

	var UkrNewName = "Ukraine Updated"
	ukr, err = s.app.UpdateCountry(ctx, ukr.ID, app.UpdateCountryParams{
		Name: &UkrNewName,
	})
	if err != nil {
		t.Fatalf("UpdateCountry: %v", err)
	}

	ukr2, err := s.app.GetCountryByID(ctx, ukr.ID)
	if err != nil {
		t.Fatalf("GetCountry: %v", err)
	}
	if ukr2.Name != UkrNewName {
		t.Errorf("expected country name '%s', got '%s'", UkrNewName, ukr2.Name)
	}

	usa, err := s.app.CreateCountry(ctx, "USA")
	if err != nil {
		t.Fatalf("CreateCountry: %v", err)
	}

	if usa.Name != "USA" {
		t.Errorf("expected country name 'USA', got '%s'", usa.Name)
	}

	countries, pag, err := s.app.ListCountries(ctx, app.FilterCountriesListParams{}, pagi.Request{
		Page: 1,
		Size: 10,
	}, []pagi.SortField{
		{Field: "name", Ascend: true},
	})
	if err != nil {
		t.Fatalf("ListCountries: %v", err)
	}

	if len(countries) != 2 {
		t.Errorf("expected 1 country, got %d", len(countries))
	}

	if pag.Total != 2 {
		t.Errorf("expected total 1 country, got %d", pag.Total)
	}

	if countries[0].ID != ukr.ID {
		t.Errorf("expected first country ID '%s', got '%s'", ukr.ID, countries[0].ID)
	}

	if ukr.Status != enum.CountryStatusUnsupported {
		t.Errorf("expected country status 'unsuported', got '%s'", ukr.Status)
	}

	ukr, err = s.app.SetCountryStatusDeprecated(ctx, ukr.ID)
	if err != nil {
		t.Fatalf("SetCountryStatusDeprecated: %v", err)
	}
	if ukr.Status != enum.CountryStatusDeprecated {
		t.Errorf("expected country status 'deprecated', got '%s'", ukr.Status)
	}

	ukr, err = s.app.SetCountryStatusSupported(ctx, ukr.ID)
	if err != nil {
		t.Fatalf("SetCountryStatusSupported: %v", err)
	}
	if ukr.Status != enum.CountryStatusSupported {
		t.Errorf("expected country status 'supported', got '%s'", ukr.Status)
	}

	usa, err = s.app.SetCountryStatusDeprecated(ctx, usa.ID)
	if err != nil {
		t.Fatalf("SetCountryStatusDeprecated: %v", err)
	}
	if usa.Status != enum.CountryStatusDeprecated {
		t.Errorf("expected country status 'deprecated', got '%s'", usa.Status)
	}

	usa, err = s.app.SetCountryStatusSupported(ctx, usa.ID)
	if err != nil {
		t.Fatalf("SetCountryStatusSupported: %v", err)
	}
	if usa.Status != enum.CountryStatusSupported {
		t.Errorf("expected country status 'supported', got '%s'", usa.Status)
	}
}
