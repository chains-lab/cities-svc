package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

func TestMayorCreateInviteNotOfficialCity(t *testing.T) {
	s, err := newSetup(t)
	cleanDb(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	ctx := context.Background()

	ukr := CreateAndActivateCountry(s, t, "Ukraine")
	kyiv := CreateCity(s, t, ukr.ID, "Kyiv")

	_, err = s.app.CreateInviteMayor(ctx, kyiv.ID, uuid.New(), roles.SuperUser)
	if !errors.Is(err, errx.ErrorCannotCreateMayorInviteForNotOfficialCity) {
		t.Fatalf("expected error %v, got %v", errx.ErrorCannotCreateMayorInviteForNotOfficialCity, err)
	}
}

func TestMayorCreateInvite(t *testing.T) {
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

	inv, err := s.app.CreateInviteMayor(ctx, kyiv.ID, uuid.New(), roles.SuperUser)
	if err != nil {
		t.Fatalf("CreateInviteMayor: %v", err)
	}

	userID := uuid.New()

	inv, err = s.app.AnswerToInvite(ctx, userID, inv.Token, enum.InviteStatusAccepted)
	if err != nil {
		t.Fatalf("AnswerToInvite: %v", err)
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

	kyivMayor, err := s.app.GetCityMayor(ctx, kyiv.ID)
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

	_, err = s.app.AnswerToInvite(ctx, secondUserID, inv.Token, enum.InviteStatusAccepted)
	if !errors.Is(err, errx.ErrorInviteAlreadyAnswered) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInviteAlreadyAnswered, err)
	}
}

func TestMayorCreateInviteDecline(t *testing.T) {
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

	inv, err := s.app.CreateInviteMayor(ctx, kyiv.ID, uuid.New(), roles.SuperUser)
	if err != nil {
		t.Fatalf("CreateInviteMayor: %v", err)
	}

	userID := uuid.New()

	inv, err = s.app.AnswerToInvite(ctx, userID, inv.Token, enum.InviteStatusRejected)
	if err != nil {
		t.Fatalf("AnswerToInvite: %v", err)
	}
	if inv.Status != enum.InviteStatusRejected {
		t.Errorf("expected invite status 'rejected', got '%s'", inv.Status)
	}
	if inv.AnsweredAt == nil {
		t.Errorf("expected invite AnsweredAt to be set, got nil")
	}
	if inv.UserID != nil && *inv.UserID != userID {
		t.Errorf("expected invite UserID to be set to %s, got %s", userID, inv.UserID)
	}

	_, err = s.app.GetCityMayor(ctx, kyiv.ID)
	if !errors.Is(err, errx.ErrorCityGovNotFound) {
		t.Fatalf("expected error %v, got %v", errx.ErrorCityGovNotFound, err)
	}

	secondUserID := uuid.New()
	_, err = s.app.AnswerToInvite(ctx, secondUserID, inv.Token, enum.InviteStatusAccepted)
	if !errors.Is(err, errx.ErrorInviteAlreadyAnswered) {
		t.Fatalf("expected error %v, got %v", errx.ErrorInviteAlreadyAnswered, err)
	}
}
