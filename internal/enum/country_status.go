package enum

type CountryStatus string

const (
	CountryStatusSupported   CountryStatus = "supported"
	CountryStatusSuspended   CountryStatus = "suspended"
	CountryStatusUnsupported CountryStatus = "unsupported"
)

var countryStatuses = []CountryStatus{
	CountryStatusSupported,
	CountryStatusSuspended,
	CountryStatusUnsupported,
}

func ParseCountryStatus(status string) (CountryStatus, bool) {
	for _, s := range countryStatuses {
		if s == CountryStatus(status) {
			return s, true
		}
	}

	return "", false
}

func GetAllCountriesStatuses() []CountryStatus {
	return countryStatuses
}
