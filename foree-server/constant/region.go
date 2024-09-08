package constant

type Region struct {
	Code    string
	IsoCode string
	Country string
	Name    string
}

var Regions = map[string]map[string]Region{
	"PK": {
		"PK-JK": {
			Code:    "JK",
			IsoCode: "PK-JK",
			Name:    "Azad Jammu and Kashmir",
			Country: "PK",
		},
		"PK-BA": {
			Code:    "BA",
			IsoCode: "PK-BA",
			Name:    "Balochistan",
			Country: "PK",
		},
		"PK-GB": {
			Code:    "GB",
			IsoCode: "PK-GB",
			Name:    "Gilgit-Baltistan",
			Country: "PK",
		},
		"PK-IS": {
			Code:    "IS",
			IsoCode: "PK-IS",
			Name:    "Islamabad",
			Country: "PK",
		},
		"PK-KP": {
			Code:    "KP",
			IsoCode: "PK-KP",
			Name:    "Khyber Pakhtunkhwa",
			Country: "PK",
		},
		"PK-PB": {
			Code:    "PB",
			IsoCode: "PK-PB",
			Name:    "Punjab",
			Country: "PK",
		},
		"PK-SD": {
			Code:    "SD",
			IsoCode: "PK-SD",
			Name:    "Sindh",
			Country: "PK",
		},
	},
	"CA": {
		"CA-AB": {
			Code:    "AB",
			IsoCode: "CA-AB",
			Name:    "Alberta",
			Country: "CA",
		},
		"CA-BC": {
			Code:    "BC",
			IsoCode: "CA-BC",
			Name:    "British Columbia",
			Country: "CA",
		},
		"CA-MB": {
			Code:    "MB",
			IsoCode: "CA-MB",
			Name:    "Manitoba",
			Country: "CA",
		},
		"CA-NB": {
			Code:    "NB",
			IsoCode: "CA-NB",
			Name:    "New Brunswick",
			Country: "CA",
		},
		"CA-NL": {
			Code:    "NL",
			IsoCode: "CA-NL",
			Name:    "Newfoundland and Labrador",
			Country: "CA",
		},
		"CA-NT": {
			Code:    "NT",
			IsoCode: "CA-NT",
			Name:    "Northwest Territories",
			Country: "CA",
		},
		"CA-NS": {
			Code:    "NS",
			IsoCode: "CA-NS",
			Name:    "Nova Scotia",
			Country: "CA",
		},
		"CA-NU": {
			Code:    "NU",
			IsoCode: "CA-NU",
			Name:    "Nunavut",
			Country: "CA",
		},
		"CA-ON": {
			Code:    "ON",
			IsoCode: "CA-ON",
			Name:    "Ontario",
			Country: "CA",
		},
		"CA-PE": {
			Code:    "PE",
			IsoCode: "CA-PE",
			Name:    "Prince Edward Island",
			Country: "CA",
		},
		"CA-QC": {
			Code:    "QC",
			IsoCode: "CA-QC",
			Name:    "Quebec",
			Country: "CA",
		},
		"CA-SK": {
			Code:    "SK",
			IsoCode: "CA-SK",
			Name:    "Saskatchewan",
			Country: "CA",
		},
		"CA-YT": {
			Code:    "YT",
			IsoCode: "CA-YT",
			Name:    "Yukon",
			Country: "CA",
		},
	},
}
