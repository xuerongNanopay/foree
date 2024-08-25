package service

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var phoneNumberReplayer = strings.NewReplacer(" ", "", "(", "", ")", "", "-", "", "+", "")
var validate = validator.New()
