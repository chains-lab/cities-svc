package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/cities-svc/internal/events/contracts"
)

type CityCreatedData struct {
	City models.City `json:"city"`
}

const CityCreatedEvent = "city.created"

func (s Service) PublishCityCreated(
	ctx context.Context,
	city models.City,
) error {
	env := contracts.Envelope[CityCreatedData]{
		Event:     CityCreatedEvent,
		Version:   "1",
		Timestamp: time.Now().UTC(),
		Data: CityCreatedData{
			City: city,
		},
	}

	return s.publish(
		ctx,
		contracts.TopicCitiesV1,
		city.ID.String(),
		env,
	)
}
