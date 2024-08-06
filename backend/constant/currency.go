package constant

var Currencies = map[string]Currency{
	"CAD": {
		Code:            "CAD",
		Num:             "124",
		Name:            "Canadian dollar",
		PostalCodeRegex: `^[A-Za-z]\d[A-Za-z][ -]?\d[A-Za-z]\d$`,
		Precision:       2,
	},
	"PKR": {
		Code:            "PKR",
		Num:             "586",
		Name:            "Pakistani rupee",
		PostalCodeRegex: `^\\d{5}$`,
		Precision:       2,
	},
	"USD": {
		Code:            "USD",
		Num:             "840",
		Name:            "United States doolar",
		PostalCodeRegex: `^\d{5}(?:[-\s]\d{4})?$`,
		Precision:       2,
	},
}

type Currency struct {
	Code            string
	Num             string
	Name            string
	PostalCodeRegex string
	Precision       int
}
