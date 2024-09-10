package constant

type Country struct {
	IsoCode         string
	IsoCode3        string
	Num             string
	Name            string
	PhoneRegex      string
	PostalCodeRegex string
}

var Countires = map[string]Country{
	"CA": {
		IsoCode:         "CA",
		IsoCode3:        "CAN",
		Num:             "124",
		Name:            "Canada",
		PhoneRegex:      "^(\\+\\d{1,2}\\s?)?1?\\-?\\.?\\s?\\(?\\d{3}\\)?[\\s.-]?\\d{3}[\\s.-]?\\d{4}$",
		PostalCodeRegex: "^[A-Za-z]\\d[A-Za-z][ -]?\\d[A-Za-z]\\d$",
	},
	"PK": {
		IsoCode:         "PK",
		IsoCode3:        "PAK",
		Num:             "586",
		Name:            "Pakistan",
		PhoneRegex:      "^((\\+92)?(0092)?(92)?(0)?)(3)([0-9]{9})$/gm",
		PostalCodeRegex: "^\\d{5}$",
	},
	"US": {
		IsoCode:         "US",
		IsoCode3:        "USA",
		Num:             "840",
		Name:            "United States of America",
		PhoneRegex:      "^(\\+\\d{1,2}\\s?)?1?\\-?\\.?\\s?\\(?\\d{3}\\)?[\\s.-]?\\d{3}[\\s.-]?\\d{4}$",
		PostalCodeRegex: "^\\d{5}(?:[-\\s]\\d{4})?$",
	},
}
