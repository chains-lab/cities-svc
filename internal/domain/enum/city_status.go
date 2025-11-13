package enum

import "fmt"

const (
	CityStatusSupported   = "supported"
	CityStatusSuspended   = "suspended"
	CityStatusUnsupported = "unsupported"
)

var cityStatuses = []string{
	CityStatusSupported,
	CityStatusSuspended,
	CityStatusUnsupported,
}

var ErrorInvalidCityStatus = fmt.Errorf("invalid city status must be one of: %s", GetAllCityStatuses())

func CheckCityStatus(status string) error {
	for _, s := range cityStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", status, ErrorInvalidCityStatus)
}

func GetAllCityStatuses() []string {
	return cityStatuses
}
