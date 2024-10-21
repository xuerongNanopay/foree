package validation

import "strings"

var PhoneNumberReplayer = strings.NewReplacer(" ", "", "(", "", ")", "", "-", "", "+", "")
