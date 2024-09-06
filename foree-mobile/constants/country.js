const Countries = {
	"CA": {
		isoCode:  "CA",
		isoCode3: "CAN",
		Num:   "124",
		Name:  "Canada",
    phoneRegex: "^(\\+\\d{1,2}\\s?)?1?\\-?\\.?\\s?\\(?\\d{3}\\)?[\\s.-]?\\d{3}[\\s.-]?\\d{4}$",
    postalCodeRegex: "^[A-Za-z]\\d[A-Za-z][ -]?\\d[A-Za-z]\\d$"
	},
	"PK": {
		isoCode:  "PK",
		isoCode3: "PAK",
		Num:   "586",
		Name:  "Pakistan",
    phoneRegex: "^((\\+92)?(0092)?(92)?(0)?)(3)([0-9]{9})$/gm",
    postalCodeRegex: "^\\d{5}$"
	},
	"US": {
		isoCode:  "US",
		isoCode3: "USA",
		Num:   "840",
		Name:  "United States of America",
    phoneRegex: "^(\\+\\d{1,2}\\s?)?1?\\-?\\.?\\s?\\(?\\d{3}\\)?[\\s.-]?\\d{3}[\\s.-]?\\d{4}$",
    postalCodeRegex: "^\\d{5}(?:[-\\s]\\d{4})?$"
	},
}

export default Countries