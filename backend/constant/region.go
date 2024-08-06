package constant

var Regions = map[string]map[string]Region{
	"PK": {
		"JK": {
			Code:    "JK",
			Code3:   "PK-JK",
			Name:    "Azad Jammu and Kashmir",
			Country: "PK",
		},
		"BA": {
			Code:    "BA",
			Code3:   "PK-BA",
			Name:    "Balochistan",
			Country: "PK",
		},
		"GB": {
			Code:    "GB",
			Code3:   "PK-GB",
			Name:    "Gilgit-Baltistan",
			Country: "PK",
		},
		"IS": {
			Code:    "IS",
			Code3:   "PK-IS",
			Name:    "Islamabad",
			Country: "PK",
		},
		"KP": {
			Code:    "KP",
			Code3:   "PK-KP",
			Name:    "Khyber Pakhtunkhwa",
			Country: "PK",
		},
		"PB": {
			Code:    "PB",
			Code3:   "PK-PB",
			Name:    "Punjab",
			Country: "PK",
		},
		"SD": {
			Code:    "SD",
			Code3:   "PK-SD",
			Name:    "Sindh",
			Country: "PK",
		},
	},
	"CA": {
		"AB": {
			Code:    "AB",
			Code3:   "CA-AB",
			Name:    "Alberta",
			Country: "CA",
		},
		"BC": {
			Code:    "BC",
			Code3:   "CA-BC",
			Name:    "British Columbia",
			Country: "CA",
		},
		"MB": {
			Code:    "MB",
			Code3:   "CA-MB",
			Name:    "Manitoba",
			Country: "CA",
		},
		"NB": {
			Code:    "NB",
			Code3:   "CA-NB",
			Name:    "New Brunswick",
			Country: "CA",
		},
		"NL": {
			Code:    "NL",
			Code3:   "CA-NL",
			Name:    "Newfoundland and Labrador",
			Country: "CA",
		},
		"NT": {
			Code:    "NT",
			Code3:   "CA-NT",
			Name:    "Northwest Territories",
			Country: "CA",
		},
		"NS": {
			Code:    "NS",
			Code3:   "CA-NS",
			Name:    "Nova Scotia",
			Country: "CA",
		},
		"NU": {
			Code:    "NU",
			Code3:   "CA-NU",
			Name:    "Nunavut",
			Country: "CA",
		},
		"ON": {
			Code:    "ON",
			Code3:   "CA-ON",
			Name:    "Ontario",
			Country: "CA",
		},
		"PE": {
			Code:    "PE",
			Code3:   "CA-PE",
			Name:    "Prince Edward Island",
			Country: "CA",
		},
		"QC": {
			Code:    "QC",
			Code3:   "CA-QC",
			Name:    "Quebec",
			Country: "CA",
		},
		"SK": {
			Code:    "SK",
			Code3:   "CA-SK",
			Name:    "Saskatchewan",
			Country: "CA",
		},
		"YT": {
			Code:    "YT",
			Code3:   "CA-YT",
			Name:    "Yukon",
			Country: "CA",
		},
	},
}

type Region struct {
	Code    string
	Code3   string
	Country string
	Name    string
}
