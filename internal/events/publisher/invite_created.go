package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	events "github.com/chains-lab/cities-svc/internal/events/contracts"
	"github.com/google/uuid"
)

type InviteCreatedData struct {
	Invite     models.Invite      `json:"invite"`
	City       models.City        `json:"city"`
	Recipients *PayloadRecipients `json:"recipients,omitempty"`
}

const InviteCreatedEvent = "city.invite.create"

func (s Service) PublishInviteCreated(
	ctx context.Context,
	invite models.Invite,
	city models.City,
	recipients ...uuid.UUID,
) error {
	event := events.Envelope[InviteCreatedData]{
		Event:     InviteCreatedEvent,
		Version:   "1",
		Timestamp: time.Now().UTC(),
		Data: InviteCreatedData{
			Invite: invite,
			City:   city,
		},
	}
	if len(recipients) > 0 {
		event.Data.Recipients = &PayloadRecipients{
			Users: recipients,
		}
	}

	return s.publish(
		ctx,
		events.TopicCitiesAdminV1,
		invite.ID.String(),
		event,
	)
}
