package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/citymod"
	"github.com/chains-lab/cities-svc/test"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func CreateMayor(s Setup, t *testing.T, cityID, userID uuid.UUID) models.CityModer {
	_, tkn, err := s.domain.moder.CreateInvite(
		context.Background(),
		enum.CityGovRoleMayor,
		cityID,
		time.Hour*24,
	)
	if err != nil {
		t.Fatalf("CreateInviteMayor: %v", err)
	}

	mod, err := s.domain.moder.AcceptInvite(context.Background(), userID, string(tkn))
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
	}

	return mod
}

func TestGetModerator(t *testing.T) {
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

	mayorID := uuid.New()
	mayor := CreateMayor(s, t, kyiv.ID, mayorID)

	gotMayor, err := s.domain.moder.Get(ctx, citymod.GetFilters{
		UserID: &mayorID,
	})
	if err != nil {
		t.Fatalf("GetMayor: %v", err)
	}
	if gotMayor.UserID != mayor.UserID {
		t.Errorf("expected mayor ID to be %s, got %s", mayor.UserID, gotMayor.UserID)
	}

	_, tkn, err := s.domain.moder.CreateInvite(
		context.Background(),
		enum.CityGovRoleModerator,
		kyiv.ID,
		time.Hour*24,
	)
	if err != nil {
		t.Fatalf("CreateInviteMayor: %v", err)
	}

	moderID := uuid.New()

	moder, err := s.domain.moder.AcceptInvite(context.Background(), moderID, string(tkn))
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
	}
	gotModer, err := s.domain.moder.Get(ctx, citymod.GetFilters{
		UserID: &moderID,
	})
	if err != nil {
		t.Fatalf("GetModerator: %v", err)
	}
	if gotModer.UserID != moder.UserID {
		t.Errorf("expected moderator ID to be %s, got %s", moder.UserID, gotModer.UserID)
	}
}

func TestCreateInviteMayor(t *testing.T) {
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

	InviteForAdmin, tkn, err := s.domain.moder.CreateInvite(ctx, enum.CityGovRoleMayor, kyiv.ID, time.Hour)
	if err != nil {
		t.Fatalf("CreateInvite: %v", err)
	}

	adminID := uuid.New()
	admin, err := s.domain.moder.AcceptInvite(ctx, adminID, string(tkn))
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
	}

	if admin.Role != enum.CityGovRoleMayor {
		t.Errorf("expected city moderator role to be 'mayor', got '%s'", admin.Role)
	}
	if admin.UserID != adminID {
		t.Errorf("expected city moderator ID to be %s, got %s", adminID, admin.UserID)
	}

	InviteForAdmin, err = s.domain.moder.GetInvite(ctx, InviteForAdmin.ID)
	if err != nil {
		t.Fatalf("GetInvite: %v", err)
	}
	if InviteForAdmin.Status != enum.InviteStatusAccepted {
		t.Errorf("expected invite status 'accepted', got '%s'", InviteForAdmin.Status)
	}
	if InviteForAdmin.AnsweredAt == nil {
		t.Errorf("expected invite AnsweredAt to be set, got nil")
	}
	if InviteForAdmin.UserID == nil || *InviteForAdmin.UserID != adminID {
		t.Errorf("expected invite UserID to be '%s', got '%v'", adminID, InviteForAdmin.UserID)
	}

	userModerator := uuid.New()

	InviteForModer, tkn, err := s.domain.moder.CreateInvite(ctx, enum.CityGovRoleMayor, kyiv.ID, time.Hour)
	if err != nil {
		t.Fatalf("CreateInvite: %v", err)
	}

	moder, err := s.domain.moder.AcceptInvite(ctx, userModerator, string(tkn))
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
	}

	InviteForModer, err = s.domain.moder.GetInvite(ctx, InviteForModer.ID)
	if err != nil {
		t.Fatalf("GetInvite: %v", err)
	}
	if InviteForModer.Status != enum.InviteStatusAccepted {
		t.Errorf("expected invite status 'accepted', got '%s'", InviteForModer.Status)
	}
	if InviteForModer.AnsweredAt == nil {
		t.Errorf("expected invite AnsweredAt to be set, got nil")
	}
	if InviteForModer.UserID == nil || *InviteForModer.UserID != userModerator {
		t.Errorf("expected invite UserID to be '%s', got '%v'", userModerator, InviteForModer.UserID)
	}

	if moder.Role != enum.CityGovRoleMayor {
		t.Errorf("expected city moderator role to be 'mayor', got '%s'", moder.Role)
	}
	if moder.UserID != userModerator {
		t.Errorf("expected city moderator ID to be %s, got %s", userModerator, moder.UserID)
	}
}

func TestTransferMayor(t *testing.T) {
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

	oldMayorId := uuid.New()
	mayor := CreateMayor(s, t, kyiv.ID, oldMayorId)

	newMayorUserID := uuid.New()

	inv, tkn, err := s.domain.moder.CreateInvite(ctx, enum.CityGovRoleMayor, kyiv.ID, time.Hour*24)
	if err != nil {
		t.Fatalf("CreateInviteMayor: %v", err)
	}

	newMayor, err := s.domain.moder.AcceptInvite(ctx, newMayorUserID, string(tkn))
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
	}
	if newMayor.Role != enum.CityGovRoleMayor {
		t.Errorf("expected new mayor role to be, got '%s'", newMayor.Role)
	}
	if newMayor.UserID != newMayorUserID {
		t.Errorf("expected new mayor ID to be %s, got %s", newMayorUserID, newMayor.UserID)
	}

	inv, err = s.domain.moder.GetInvite(ctx, inv.ID)
	if err != nil {
		t.Fatalf("GetInvite: %v", err)
	}
	if inv.Status != enum.InviteStatusAccepted {
		t.Errorf("expected invite status 'sent', got '%s'", inv.Status)
	}
	if inv.AnsweredAt == nil {
		t.Errorf("expected invite AnsweredAt to be set, got nil")
	}

	mayor, err = s.domain.moder.Get(context.Background(), citymod.GetFilters{
		CityID: &kyiv.ID,
		Role:   &mayor.Role,
	})
	if err != nil {
		t.Fatalf("GetCityMayor: %v", err)
	}
	if mayor.UserID != newMayorUserID {
		t.Errorf("expected city mayor ID to be %s, got %s", newMayorUserID, mayor.UserID)
	}
	if mayor.Role != enum.CityGovRoleMayor {
		t.Errorf("expected city mayor role to be 'mayor', got '%s'", mayor.Role)
	}

	_, err = s.domain.moder.Get(ctx, citymod.GetFilters{
		UserID: &oldMayorId,
	})
	if !errors.Is(err, errx.ErrorCityGovNotFound) {
		t.Fatalf("expected error %v, got %v", errx.ErrorCityGovNotFound, err)
	}
}

func TestMayorCreateInviteNotOfficialCity(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	_, _, err = s.domain.moder.CreateInvite(ctx, enum.CityGovRoleMayor, kyiv.ID, time.Hour*24)
	if !errors.Is(err, errx.ErrorCityIsNotSupported) {
		t.Fatalf("expected error %v, got %v", errx.ErrorCityIsNotSupported, err)
	}
}

func TestMayorCreateInvite(t *testing.T) {
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

	inv, tkn, err := s.domain.moder.CreateInvite(ctx, enum.CityGovRoleMayor, kyiv.ID, time.Hour*24)
	if err != nil {
		t.Fatalf("CreateInviteMayor: %v", err)
	}

	userID := uuid.New()

	_, err = s.domain.moder.AcceptInvite(ctx, userID, string(tkn))
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
	}

	inv, err = s.domain.moder.GetInvite(ctx, inv.ID)
	if err != nil {
		t.Fatalf("GetInvite: %v", err)
	}
	if inv.Status != enum.InviteStatusAccepted {
		t.Errorf("expected invite status 'accepted', got '%s'", inv.Status)
	}
	if inv.AnsweredAt == nil {
		t.Errorf("expected invite AnsweredAt to be set, got nil")
	}
	if inv.UserID != nil && *inv.UserID != userID {
		t.Errorf("expected invite UserID to be set to %s, got %s", userID, inv.UserID)
	}

	mayorR := enum.CityGovRoleMayor
	kyivMayor, err := s.domain.moder.Get(ctx, citymod.GetFilters{
		CityID: &kyiv.ID,
		Role:   &mayorR,
	})
	if err != nil {
		t.Fatalf("GetCityMayor: %v", err)
	}
	if kyivMayor.UserID != userID {
		t.Errorf("expected city mayor ID to be %s, got %s", userID, kyivMayor.UserID)
	}
	if kyivMayor.Role != enum.CityGovRoleMayor {
		t.Errorf("expected city mayor role to be 'mayor', got '%s'", kyivMayor.Role)
	}

	secondUserID := uuid.New()

	_, err = s.domain.moder.AcceptInvite(ctx, secondUserID, string(tkn))
	if !errors.Is(err, errx.ErrorInviteAlreadyAnswered) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInviteAlreadyAnswered, err)
	}
}
