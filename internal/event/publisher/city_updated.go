package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	events "github.com/chains-lab/cities-svc/internal/event/contracts"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type CityUpdatedData struct {
	ID        uuid.UUID `json:"id"`
	CountryID string    `json:"country_id"`
	Point     orb.Point `json:"point"`
	Status    string    `json:"status"`
	Name      string    `json:"name"`
	Icon      *string   `json:"icon"`
	Slug      *string   `json:"slug"`
	Timezone  string    `json:"timezone"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const CityUpdateEvent = "city.admin.update"

func (s Service) PublishCityUpdated(
	ctx context.Context,
	city models.City,
) error {
	env := events.Envelope[CityUpdatedData]{
		Event:     CityUpdateEvent,
		Version:   "1",
		Timestamp: time.Now().UTC(),
		Data: CityUpdatedData{
			ID:        city.ID,
			CountryID: city.CountryID,
			Point:     city.Point,
			Status:    city.Status,
			Name:      city.Name,
			Icon:      city.Icon,
			Slug:      city.Slug,
			Timezone:  city.Timezone,
			CreatedAt: city.CreatedAt,
			UpdatedAt: city.UpdatedAt,
		},
	}

	return s.publish(
		ctx,
		events.TopicCitiesAdminV1,
		city.ID.String(),
		env,
	)
}
