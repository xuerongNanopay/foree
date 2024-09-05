package constant

var Countires = map[string]Country{
	"CA": {
		Code:  "CA",
		Code3: "CAN",
		Num:   "124",
		Name:  "Canada",
	},
	"PK": {
		Code:  "PK",
		Code3: "PAK",
		Num:   "586",
		Name:  "Pakistan",
	},
	"US": {
		Code:  "US",
		Code3: "USA",
		Num:   "840",
		Name:  "United States of America",
	},
}

type Country struct {
	Code  string
	Code3 string
	Num   string
	Name  string
}
