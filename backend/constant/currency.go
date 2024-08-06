package constant

var Currencies = map[string]Currency{
	"CAD": {
		Code:      "CAD",
		Num:       "124",
		Name:      "Canadian dollar",
		Precision: 2,
	},
	"PKR": {
		Code:      "PKR",
		Num:       "586",
		Name:      "Pakistani rupee",
		Precision: 2,
	},
	"USD": {
		Code:      "USD",
		Num:       "840",
		Name:      "United States doolar",
		Precision: 2,
	},
}

type Currency struct {
	Code      string
	Num       string
	Name      string
	Precision int
}
