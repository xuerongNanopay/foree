package constant

type Country struct {
	IsoCode  string
	IsoCode3 string
	Num      string
	Name     string
}

var Countires = map[string]Country{
	"CA": {
		IsoCode:  "CA",
		IsoCode3: "CAN",
		Num:      "124",
		Name:     "Canada",
	},
	"PK": {
		IsoCode:  "PK",
		IsoCode3: "PAK",
		Num:      "586",
		Name:     "Pakistan",
	},
	"US": {
		IsoCode:  "US",
		IsoCode3: "USA",
		Num:      "840",
		Name:     "United States of America",
	},
}
