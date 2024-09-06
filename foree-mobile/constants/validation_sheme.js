import { number, string } from 'yup'
import Countries from './country'

const internationalNameRegex = /^[a-zA-ZàáâäãåąčćęèéêëėįìíîïłńòóôöõøùúûüųūÿýżźñçčšžæÀÁÂÄÃÅĄĆČĖĘÈÉÊËÌÍÎÏĮŁŃÒÓÔÖÕØÙÚÛÜŲŪŸÝŻŹÑßÇŒÆČŠŽ∂ð ,.'-]+$/u

//Excape empty string check
const emptyRegexWrapper = (regex) => `(${regex})|(^$)`
const String = ({
  min=0,
  max=256,
  required=true
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  return ret.min(min, min!==max ? `at least ${min} characters`: `must be ${min} characters`).max(max, min!==max ? `at most ${max} characters`: `must be ${max} characters`)
}

const NameScheme = ({
  min=0,
  max=256,
  required=true
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  return ret.matches(new RegExp(emptyRegexWrapper("^[A-Za-z ,.'-]+$")), "invalid character").min(min, min!==max ? `at least ${min} characters`: `must be ${min} characters`).max(max, min!==max ? `at most ${max} characters`: `must be ${max} characters`)
}

const AlphaNumberScheme = ({
  min=0,
  max=256,
  required=true
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  return ret.alphanum().min(min, min!==max ? `at least ${min} characters`: `must be ${min} characters`).max(max, min!==max ? `at most ${max} characters`: `must be ${max} characters`)
}

const EmailScheme = ({
  required=true
})=> {
  const ret = required ? string().trim().required("required") : string().trim()
  return ret.email("invalid email")
}

const NumberOnlyScheme = ({
  min=0,
  max=256, 
  required=true
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  return ret.matches(new RegExp(emptyRegexWrapper("^\\d+$")), "").min(min, min!==max ? `at least ${min} characters`: `must be ${min} characters`).max(max, min!==max ? `at most ${max} characters`: `must be ${max} characters`)
}

const AlphaOnlyScheme = ({
  min=0,
  max=256,
  upperCaseOnly=false,
  required=true
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  const regex = upperCaseOnly ? new RegExp(emptyRegexWrapper("[A-Z]+$")) : new RegExp(emptyRegexWrapper("^[A-Za-z]+$"))
  return ret.matches(regex, "invalid character").min(min, min!==max ? `at least ${min} characters`: `must be ${min} characters`).max(max, min!==max ? `at most ${max} characters`: `must be ${max} characters`)
}

const DateOnlyScheme =({
  required=true
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  return ret.matches(new RegExp(emptyRegexWrapper("^\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])$")), "must be YYYY-MM-DD")
}

const IntegerScheme = ({
  min,
  max, 
  required=true
}) => {
  let ret = required ? number().integer("must be integer").required('required') : number().integer()
  ret = typeof(min) != "undefined" ? ret.min(min, `at least ${min}`) : ret
  ret = typeof(max) != "undefined" ? ret.min(max, `at most ${max}`) : ret
  return ret
}

const PositiveIntegerScheme = ({
  max, 
  required=true
}) => {
  return IntegerScheme({min: 0, max, required})
}

const FloatScheme = ({
  min,
  max, 
  required=true
}) => {
  let ret = required ? number().required('required') : number()
  ret = typeof(min) != "undefined" ? ret.min(min, `at least ${min}`) : ret
  ret = typeof(max) != "undefined" ? ret.min(max, `at most ${max}`) : ret
  return ret
}

const PositiveFloatScheme = ({
  max, 
  required=true
}) => {
  return FloatScheme({min: 0, max, required})
}

// Minimum eight characters, at least one letter and one number:
const PasswdMinFourOneLetterOneNumber = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{4,}$/
// Minimum eight characters, at least one letter and one number:
const PasswdMinEightOneLetterOneNumber = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/
// Minimum eight characters, at least one letter, one number and one special character:
const PasswdMinEightOneLetterOneNumberOneSpecial = /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/
// Minimum eight characters, at least one uppercase letter, one lowercase letter and one number:
const PasswdMinEightOneUpperOneLowerOneNumber = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/
// Minimum eight characters, at least one uppercase letter, one lowercase letter, one number and one special character
const PasswdMinEightOneUpperOneLowerOneNumberOneSepcia = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/

const PasswordScheme = ({
  passwdLevel=PasswdMinFourOneLetterOneNumber
}) => {
  return string().trim().required("required").matches(passwdLevel)
}

const PostalCodeScheme = ({
  required=true,
  countryCode,
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  return ret.matches(new RegExp(emptyRegexWrapper(Countries[countryCode].postalCodeRegex)), "invalid postal code")
}

const PhoneNumber = ({
  required=true,
  countryCode,
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  return  ret.matches(new RegExp(emptyRegexWrapper(Countries[countryCode].phoneRegex)), "invalid phone number")
}

export default {
  String,
  NameScheme,
  AlphaNumberScheme,
  EmailScheme,
  NumberOnlyScheme,
  DateOnlyScheme,
  AlphaOnlyScheme,
  IntegerScheme,
  PositiveIntegerScheme,
  FloatScheme,
  PositiveFloatScheme,
  PasswordScheme,
  PostalCodeScheme,
  PhoneNumber,
  PasswdMinFourOneLetterOneNumber,
  PasswdMinEightOneLetterOneNumber,
  PasswdMinEightOneLetterOneNumberOneSpecial,
  PasswdMinEightOneUpperOneLowerOneNumber,
  PasswdMinEightOneUpperOneLowerOneNumberOneSepcia
}
