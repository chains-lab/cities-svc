package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/errx"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/cities-svc/test"
	"github.com/google/uuid"
)

func CreateHead(s Setup, t *testing.T, cityID, userID uuid.UUID) models.CityAdmin {
	inv, err := s.domain.moder.CreateInvite(
		context.Background(),
		enum.CityAdminRoleHead,
		cityID,
		time.Hour*24,
	)
	if err != nil {
		t.Fatalf("CreateInviteMayor: %v", err)
	}

	inv, err = s.domain.moder.AcceptInvite(context.Background(), userID, inv.Token)
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
	}

	mod, err := s.domain.moder.Get(context.Background(), admin.GetFilters{
		UserID: &userID,
	})
	if err != nil {
		t.Fatalf("GetHead: %v", err)
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

	headID := uuid.New()
	head := CreateHead(s, t, kyiv.ID, headID)

	gotMayor, err := s.domain.moder.Get(ctx, admin.GetFilters{
		UserID: &headID,
	})
	if err != nil {
		t.Fatalf("GetMayor: %v", err)
	}
	if gotMayor.UserID != head.UserID {
		t.Errorf("expected head ID to be %s, got %s", head.UserID, gotMayor.UserID)
	}

	inv, err := s.domain.moder.CreateInvite(
		context.Background(),
		enum.CityAdminRoleModerator,
		kyiv.ID,
		time.Hour*24,
	)
	if err != nil {
		t.Fatalf("CreateInviteMayor: %v", err)
	}

	moderID := uuid.New()

	inv, err = s.domain.moder.AcceptInvite(context.Background(), moderID, inv.Token)
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
	}

	gotModer, err := s.domain.moder.Get(ctx, admin.GetFilters{
		UserID: &moderID,
	})
	if err != nil {
		t.Fatalf("GetModerator: %v", err)
	}
	if gotModer.UserID != *inv.UserID {
		t.Errorf("expected moderator ID to be %s, got %s", inv.UserID, gotModer.UserID)
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

	InviteForAdmin, err := s.domain.moder.CreateInvite(ctx, enum.CityAdminRoleHead, kyiv.ID, time.Hour)
	if err != nil {
		t.Fatalf("CreateInvite: %v", err)
	}

	adminID := uuid.New()

	InviteForAdmin, err = s.domain.moder.AcceptInvite(ctx, adminID, InviteForAdmin.Token)
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
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

	InviteForModer, err := s.domain.moder.CreateInvite(ctx, enum.CityAdminRoleHead, kyiv.ID, time.Hour)
	if err != nil {
		t.Fatalf("CreateInvite: %v", err)
	}

	InviteForModer, err = s.domain.moder.AcceptInvite(ctx, userModerator, InviteForModer.Token)
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
}

func TestCreateInviteNotOfficialCity(t *testing.T) {
	s, err := newSetup(t)
	test.CleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	_, err = s.domain.moder.CreateInvite(ctx, enum.CityAdminRoleHead, kyiv.ID, time.Hour*24)
	if !errors.Is(err, errx.ErrorCityIsNotSupported) {
		t.Fatalf("expected error %v, got %v", errx.ErrorCityIsNotSupported, err)
	}
}
