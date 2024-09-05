import Joi from 'joi'

const AlphanumNumberScheme = ({
  min=8,
  max=30,
  required=true
}) => {
  const ret = Joi.string().alphanum().min(min).max(max)
  return required ? ret.required() : ret
}

const EmailScheme = ({
  registerdTLD=true,
  required=true
})=> {
  const ret = Joi.string().email({ tlds: { allow: registerdTLD } })
  return required ? ret.required() : ret
}

const NumberOnlyScheme = ({
  min=8,
  max=30, 
  required=true
}) => {
  const ret = Joi.string().regex(/^\d+$/).min(min).max(max)
  return required ? ret.required() : ret
}

const AlphaOnlyScheme = ({
  min=8,
  max=30,
  upperCaseOnly=false,
  required=true
}) => {
  const regex = upperCaseOnly ? /^[A-Z]+$/ : /^[A-Za-z]+$/
  const ret = Joi.string().regex(regex).min(min).max(max)
  return required ? ret.required() : ret
}

const DateOnlyScheme =({
  required=true
}) => {
  const ret = Joi.string().regex(/^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$/)
  return required ? ret.required() : ret
}

const IntegerScheme = ({
  min,
  max, 
  required=true
}) => {
  let ret = Joi.number().integer().max(max)
  ret = typeof(min) != "undefined" ? ret.min(min) : ret
  ret = typeof(max) != "undefined" ? ret.min(min) : ret
  return required ? ret.required() : ret
}

const PositiveIntegerScheme = ({
  max, 
  required=true
}) => {
  return IntegerScheme({min: 1, max, required})
}

const FloatScheme = ({
  min,
  max, 
  required=true
}) => {
  let ret = Joi.number().integer().max(max)
  ret = typeof(min) != "undefined" ? ret.min(min) : ret
  ret = typeof(max) != "undefined" ? ret.min(min) : ret
  return required ? ret.required() : ret
}

const PositiveFloatScheme = ({
  max, 
  required=true
}) => {
  return FloatScheme({min: 1, max, required})
}

// Minimum eight characters, at least one letter and one number:
// "^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$"

const PasswordScheme =({
  required=true
}) => {
  const ret = Joi.string().regex(/^$/)
  return required ? ret.required() : ret
}

export default {
  AlphanumNumberScheme,
  EmailScheme,
  NumberOnlyScheme,
  DateOnlyScheme,
  AlphaOnlyScheme,
  IntegerScheme,
  PositiveIntegerScheme,
  FloatScheme,
  PositiveFloatScheme
};
