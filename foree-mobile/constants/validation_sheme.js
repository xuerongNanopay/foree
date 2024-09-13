import { number, string } from 'yup'
import Countries from './country'

const internationalNameRegex = /^[a-zA-ZàáâäãåąčćęèéêëėįìíîïłńòóôöõøùúûüųūÿýżźñçčšžæÀÁÂÄÃÅĄĆČĖĘÈÉÊËÌÍÎÏĮŁŃÒÓÔÖÕØÙÚÛÜŲŪŸÝŻŹÑßÇŒÆČŠŽ∂ð ,.'-]+$/u

// TODO: min and max is not working well with not required.
//Excape empty string check
const emptyRegexWrapper = (regex) => `(${regex})|(^$)`
const String = ({min=0, max=256, required=true}={
  min,
  max,
  required
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  return ret.min(min, min!==max ? `at least ${min} characters`: `must be ${min} characters`).max(max, min!==max ? `at most ${max} characters`: `must be ${max} characters`)
}

const NameScheme = ({min=0, max=256, required=true}={
  min,
  max,
  required
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  return ret.matches(new RegExp(emptyRegexWrapper("^[A-Za-z ,.'-]+$")), "invalid character").min(min, min!==max ? `at least ${min} characters`: `must be ${min} characters`).max(max, min!==max ? `at most ${max} characters`: `must be ${max} characters`)
}

const AlphaNumberScheme = ({min=0, max=256, required=true}={
  min,
  max,
  required
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  return ret.alphanum().min(min, min!==max ? `at least ${min} characters`: `must be ${min} characters`).max(max, min!==max ? `at most ${max} characters`: `must be ${max} characters`)
}

const EmailScheme = ({required=true}={
  required
})=> {
  const ret = required ? string().trim().required("required") : string().trim()
  return ret.email("invalid email")
}

const NumberOnlyScheme = ({min=0, max=256, required=true}={
  min,
  max, 
  required
}) => {
  const ret = required ? string().trim().required("required") : string()
  return ret.matches(new RegExp(emptyRegexWrapper("^\\d+$")), "invalid digits").min(min, min!==max ? `at least ${min} characters`: `must be ${min} characters`).max(max, min!==max ? `at most ${max} characters`: `must be ${max} characters`)
}

const AlphaOnlyScheme = ({min=0, max=256, required=true, upperCaseOnly=false}={
  min,
  max,
  upperCaseOnly,
  required
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  const regex = upperCaseOnly ? new RegExp(emptyRegexWrapper("[A-Z]+$")) : new RegExp(emptyRegexWrapper("^[A-Za-z]+$"))
  return ret.matches(regex, "invalid character").min(min, min!==max ? `at least ${min} characters`: `must be ${min} characters`).max(max, min!==max ? `at most ${max} characters`: `must be ${max} characters`)
}

const CountryISOScheme = ({required=true}={
  required
}) => AlphaOnlyScheme({upperCaseOnly: true, min:2, max:2, required})

const ProvinceISOScheme = ({required=true}={
  required
}) => {
  const min = 5
  const max = 5
  const ret = required ? string().trim().required("required") : string().trim()
  const regex = new RegExp(emptyRegexWrapper("[A-Z]{2}\\-[A-Z]{2}$"))
  return ret.matches(regex, "invalid character").min(min, min!==max ? `at least ${min} characters`: `must be ${min} characters`).max(max, min!==max ? `at most ${max} characters`: `must be ${max} characters`)
}

const DateOnlyScheme =({required=true}={
  required
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  return ret.matches(new RegExp(emptyRegexWrapper("^\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])$")), "must be YYYY-MM-DD")
}

const AgeScheme =({minAge=19, required=true}={
  minAge,
  required
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  return ret
          .matches(new RegExp(emptyRegexWrapper("^\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])$")), "must be YYYY-MM-DD")
          .test("ageTest", `require ${minAge} years old`, async (value)=> {
            try {
              const birth = Date.parse(value)
              if ( !isNaN(birth) ) {
                const now = new Date()
                return (now.getTime()-birth)/(3600000*24*365) > 19
              }
            } catch (e) {
              console.log("AgeScheme", e)
            }
          })
}

const IntegerScheme = ({min, max, required=true}={
  min,
  max, 
  required
}) => {
  let ret = required ? number().integer("must be integer").required('required') : number().integer()
  ret = typeof(min) != "undefined" ? ret.min(min, `at least ${min}`) : ret
  ret = typeof(max) != "undefined" ? ret.min(max, `at most ${max}`) : ret
  return ret
}

const PositiveIntegerScheme = ({max, required=true}={
  max, 
  required
}) => {
  return IntegerScheme({min: 0, max, required})
}

const FloatScheme = ({min, max, required=true}={
  min,
  max, 
  required
}) => {
  let ret = required ? number().required('required') : number()
  ret = typeof(min) != "undefined" ? ret.min(min, `at least ${min}`) : ret
  ret = typeof(max) != "undefined" ? ret.min(max, `at most ${max}`) : ret
  return ret
}

const PositiveFloatScheme = ({max, required=true}={
  max, 
  required
}) => {
  return FloatScheme({min: 0, max, required})
}

// Minimum four characters, at least one letter and one number:
const PasswdMinFourOneLetterOneNumber = {
  regex: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{4,}$/,
  message: "minimum four characters, at least one letter and one number"
}
// Minimum eight characters, at least one letter and one number:
const PasswdMinEightOneLetterOneNumber = {
  regex: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/,
  message: "minimum eight characters, at least one letter and one number"
}
// Minimum eight characters, at least one letter, one number and one special character:
const PasswdMinEightOneLetterOneNumberOneSpecial = {
  regex: /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/,
  message: "minimum eight characters, at least one letter, one number and one special character"
}
// Minimum eight characters, at least one uppercase letter, one lowercase letter and one number:
const PasswdMinEightOneUpperOneLowerOneNumber = {
  regex: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/,
  message: "minimum eight characters, at least one uppercase letter, one lowercase letter and one number"
}
// Minimum eight characters, at least one uppercase letter, one lowercase letter, one number and one special character
const PasswdMinEightOneUpperOneLowerOneNumberOneSepcia = {
  regex: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/,
  message: "minimum eight characters, at least one uppercase letter, one lowercase letter, one number and one special character"
}

const PasswordScheme = ({passwordValidator=PasswdMinFourOneLetterOneNumber}={
  passwordValidator
}) => {
  return string().trim().required("required").matches(passwordValidator.regex, passwordValidator.message)
}

const PostalCodeScheme = ({required=true, countryCode}={
  required,
  countryCode,
}) => {
  const ret = required ? string().trim().required("required") : string().trim()
  return ret.matches(new RegExp(emptyRegexWrapper(Countries[countryCode].postalCodeRegex)), "invalid postal code")
}

const PhoneNumber = ({required=true, countryCode}={
  required,
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
  CountryISOScheme,
  PhoneNumber,
  ProvinceISOScheme,
  AgeScheme,
  PasswdMinFourOneLetterOneNumber,
  PasswdMinEightOneLetterOneNumber,
  PasswdMinEightOneLetterOneNumberOneSpecial,
  PasswdMinEightOneUpperOneLowerOneNumber,
  PasswdMinEightOneUpperOneLowerOneNumberOneSepcia
}
