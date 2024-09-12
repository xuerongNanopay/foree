const Countries = {
  "CA": {
    isoCode:  "CA",
    isoCode3: "CAN",
    isoNum:   "124",
    name:  "Canada",
    unicodeIcon: "\u{1F1E8}\u{1F1E6}",
    phoneRegex: "^(\\+\\d{1,2}\\s?)?1?\\-?\\.?\\s?\\(?\\d{3}\\)?[\\s.-]?\\d{3}[\\s.-]?\\d{4}$",
    postalCodeRegex: "^[A-Za-z]\\d[A-Za-z][ -]?\\d[A-Za-z]\\d$"
  },
  "PK": {
    isoCode:  "PK",
    isoCode3: "PAK",
    isoNum:   "586",
    name:  "Pakistan",
    unicodeIcon: "\u{1F1F5}\u{1F1F0}",
    phoneRegex: "^((\\+92)?(0092)?(92)?(0)?)(3)([0-9]{9})$/gm",
    postalCodeRegex: "^\\d{5}$"
  },
  "US": {
    isoCode:  "US",
    isoCode3: "USA",
    isoNum:   "840",
    name:  "United States of America",
    unicodeIcon: "\u{1F1FA}\u{1F1F8}",
    phoneRegex: "^(\\+\\d{1,2}\\s?)?1?\\-?\\.?\\s?\\(?\\d{3}\\)?[\\s.-]?\\d{3}[\\s.-]?\\d{4}$",
    postalCodeRegex: "^\\d{5}(?:[-\\s]\\d{4})?$"
  }
}

export default Countries