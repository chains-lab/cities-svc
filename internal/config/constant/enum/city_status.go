package enum

import "fmt"

const (
	CityStatusSupported = "supported"
	CityStatusSuspended = "suspended"
)

var cityStatuses = []string{
	CityStatusSupported,
	CityStatusSuspended,
}

var ErrorCityStatusNotSupported = fmt.Errorf("invalid city status must be one of: %v", GetAllCitiesStatuses())

func ParseCityStatus(status string) (string, error) {
	for _, s := range cityStatuses {
		if s == status {
			return s, nil
		}
	}

	return "", fmt.Errorf("'%s', %w", status, ErrorCityStatusNotSupported)
}

func GetAllCitiesStatuses() []string {
	return cityStatuses
}
