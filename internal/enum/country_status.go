package enum

type CountryStatus int

const (
	CountryNotSupport CountryStatus = iota
	CountrySupport
)

var countriesStatus = []string{
	"not_support",
	"support",
}

func (s CountryStatus) String() string {
	if int(s) < len(countriesStatus) {
		return countriesStatus[s]
	}
	return ""
}

func CheckCountryStatus(s string) bool {
	for _, name := range countriesStatus {
		if name == s {
			return true
		}
	}
	return false
}

func GetAllCountriesStatuses() []string {
	return countriesStatus
}
