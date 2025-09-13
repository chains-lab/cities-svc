package apptest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/app/models"
	"github.com/chains-lab/cities-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
)

func CreateMayor(s Setup, t *testing.T, cityID, userID uuid.UUID) models.Gov {
	inv, err := s.app.CreateInviteMayor(context.Background(), cityID, uuid.New(), roles.SuperUser)
	if err != nil {
		t.Fatalf("CreateInviteMayor: %v", err)
	}

	inv, err = s.app.AcceptInvite(context.Background(), userID, inv.Token)
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
	}

	mayor, err := s.app.GetCityMayor(context.Background(), cityID)
	if err != nil {
		t.Fatalf("GetCityMayor: %v", err)
	}

	return mayor
}

func TestCreateInviteMayor(t *testing.T) {
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

	mayor := CreateMayor(s, t, kyiv.ID, uuid.New())

	InviteForAdmin, err := s.app.SentInvite(ctx, app.SentInviteParams{
		InitiatorID: mayor.UserID,
		CityID:      kyiv.ID,
		Role:        enum.CityGovRoleModerator,
	})
	if err != nil {
		t.Fatalf("SentInvite: %v", err)
	}

	userAdmin := uuid.New()
	InviteForAdmin, err = s.app.AcceptInvite(ctx, userAdmin, InviteForAdmin.Token)
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
	}
	if InviteForAdmin.Status != enum.InviteStatusAccepted {
		t.Errorf("expected invite status 'accepted', got '%s'", InviteForAdmin.Status)
	}

	if InviteForAdmin.AnsweredAt == nil {
		t.Errorf("expected invite AnsweredAt to be set, got nil")
	}
	if InviteForAdmin.UserID == nil || *InviteForAdmin.UserID != userAdmin {
		t.Errorf("expected invite UserID to be '%s', got '%v'", userAdmin, InviteForAdmin.UserID)
	}

	userModerator := uuid.New()

	InviteForModerator, err := s.app.SentInvite(ctx, app.SentInviteParams{
		InitiatorID: mayor.UserID,
		CityID:      kyiv.ID,
		Role:        enum.CityGovRoleModerator,
	})
	if err != nil {
		t.Fatalf("SentInvite: %v", err)
	}

	InviteForModerator, err = s.app.AcceptInvite(ctx, userModerator, InviteForModerator.Token)
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
	}
}

func TestTransferMayor(t *testing.T) {
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

	oldMayorId := uuid.New()
	mayor := CreateMayor(s, t, kyiv.ID, oldMayorId)

	newMayorUserID := uuid.New()

	inv, err := s.app.CreateInviteMayor(ctx, kyiv.ID, oldMayorId, roles.User)
	if err != nil {
		t.Fatalf("CreateInviteMayor: %v", err)
	}

	inv, err = s.app.AcceptInvite(ctx, newMayorUserID, inv.Token)
	if err != nil {
		t.Fatalf("AcceptInvite: %v", err)
	}

	mayor, err = s.app.GetCityMayor(context.Background(), inv.CityID)
	if err != nil {
		t.Fatalf("GetCityMayor: %v", err)
	}
	if mayor.UserID != newMayorUserID {
		t.Errorf("expected city mayor ID to be %s, got %s", newMayorUserID, mayor.UserID)
	}
	if mayor.Role != enum.CityGovRoleMayor {
		t.Errorf("expected city mayor role to be 'mayor', got '%s'", mayor.Role)
	}

	_, err = s.app.GetGov(ctx, oldMayorId)
	if !errors.Is(err, errx.ErrorCityGovNotFound) {
		t.Fatalf("expected error %v, got %v", errx.ErrorCityGovNotFound, err)
	}
}
