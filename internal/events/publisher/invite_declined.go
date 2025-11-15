package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/events/contracts"
	"github.com/google/uuid"
)

type InviteDeclinedData struct {
	Invite     models.Invite     `json:"invite"`
	City       models.City       `json:"city"`
	Recipients PayloadRecipients `json:"recipients"`
}

const InviteCanceledEvent = "city.invite.decline"

func (s Service) PublishInviteDeclined(
	ctx context.Context,
	invite models.Invite,
	city models.City,
	recipients ...uuid.UUID,
) error {
	return s.publish(
		ctx,
		contracts.TopicCitiesV1,
		invite.ID.String(),
		contracts.Envelope[InviteDeclinedData]{
			Event:     InviteCanceledEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: InviteDeclinedData{
				Invite: invite,
				City:   city,
				Recipients: PayloadRecipients{
					Users: recipients,
				},
			},
		},
	)
}
