package constant

import "fmt"

const (
	CityStatusOfficial   = "official"
	CityStatusCommunity  = "community"
	CityStatusArchived   = "archived"
	CityStatusDeprecated = "deprecated"
)

var cityStatuses = []string{
	CityStatusOfficial,
	CityStatusCommunity,
	CityStatusArchived,
	CityStatusDeprecated,
}

var ErrorCityStatusNotSupported = fmt.Errorf("invalid city status must be one of: %v", GetAllCitiesStatuses())

func CheckCityStatus(status string) error {
	for _, s := range cityStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", status, ErrorCityStatusNotSupported)
}

func GetAllCitiesStatuses() []string {
	return cityStatuses
}
