package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/events/contracts"
	"github.com/google/uuid"
)

type InviteAcceptedPayload struct {
	Invite     models.Invite      `json:"invite"`
	City       models.City        `json:"city"`
	CityAdmin  models.CityAdmin   `json:"city_admin"`
	Recipients *PayloadRecipients `json:"recipients,omitempty"`
}

const InviteAcceptedEvent = "city.invite.accepted"

func (s Service) PublishInviteAccepted(
	ctx context.Context,
	invite models.Invite,
	city models.City,
	cityAdmin models.CityAdmin,
	recipients ...uuid.UUID,
) error {
	event := contracts.Envelope[InviteAcceptedPayload]{
		Event:     InviteAcceptedEvent,
		Version:   "1",
		Timestamp: time.Now().UTC(),
		Data: InviteAcceptedPayload{
			City:      city,
			Invite:    invite,
			CityAdmin: cityAdmin,
		},
	}
	if len(recipients) > 0 {
		event.Data.Recipients = &PayloadRecipients{
			Users: recipients,
		}
	}

	return s.publish(
		ctx,
		contracts.TopicCitiesV1,
		invite.ID.String(),
		event,
	)
}
