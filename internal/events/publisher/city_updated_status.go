package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/enum"
	"github.com/chains-lab/cities-svc/internal/domain/models"
	events "github.com/chains-lab/cities-svc/internal/events/contracts"
	"github.com/google/uuid"
)

type UpdatedStatusStatusData struct {
	City       models.City       `json:"city"`
	Recipients PayloadRecipients `json:"recipients"`
}

const CityUpdatedStatusSupportedEvent = "city.update.status.supported"
const CityUpdatedStatusSuspendedEvent = "city.update.status.suspended"
const CityUpdatedStatusUnsupportedEvent = "city.update.status.unsupported"

func (s Service) PublishCityUpdatedStatus(
	ctx context.Context,
	city models.City,
	status string,
	recipients []uuid.UUID,
) error {
	var eventName string
	switch status {
	case enum.CityStatusSupported:
		eventName = CityUpdatedStatusSupportedEvent
	case enum.CityStatusSuspended:
		eventName = CityUpdatedStatusSuspendedEvent
	case enum.CityStatusUnsupported:
		eventName = CityUpdatedStatusUnsupportedEvent
	default:
		return enum.ErrorInvalidCityStatus
	}

	return s.publish(
		ctx,
		events.TopicCitiesAdminV1,
		city.ID.String(),
		events.Envelope[CityUpdatedData]{
			Event:     eventName,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: CityUpdatedData{
				City: city,
				Recipients: PayloadRecipients{
					Users: recipients,
				},
			},
		},
	)
}
