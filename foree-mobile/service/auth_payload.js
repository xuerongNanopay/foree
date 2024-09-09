import { object, string } from 'yup'
import fieldScheme from '../constants/validation_sheme'

const SignUpScheme = object({
  email: fieldScheme.EmailScheme({}),
  password: fieldScheme.PasswordScheme({}),
})

const VerifyEmailScheme = object({
  code: fieldScheme.NumberOnlyScheme({min:6, max:6}),
})

const ChangePasswdScheme = object({
  oldPassword: string().required(),
  password: fieldScheme.PasswordScheme({}),
})

const LoginScheme = object({
  email: fieldScheme.EmailScheme({}),
  password: fieldScheme.PasswordScheme({}),
})

const ForgetPasswdScheme = object({
  email: fieldScheme.EmailScheme({}),
})

const ForgetPasswdUpdateScheme = object({
  retrieveCode: string().required(),
  newPassword: fieldScheme.PasswordScheme({}),
})

const OnboardingScheme = object({
  firstName: fieldScheme.NameScheme({}),
  middleName: fieldScheme.NameScheme({required: false}),
  lastName: fieldScheme.NameScheme({}),
  address1: fieldScheme.String({}),
  address2: fieldScheme.String({required: false}),
  city: fieldScheme.String({}),
  province: fieldScheme.ProvinceISOScheme({}),
  country: fieldScheme.CountryISOScheme({}),
  postalCode: fieldScheme.PostalCodeScheme({countryCode:"CA"}),
  phoneNumber: fieldScheme.PhoneNumber({countryCode:"CA"}),
  dob: fieldScheme.AgeScheme({}),
  pob: fieldScheme.CountryISOScheme({}),
  nationality: fieldScheme.CountryISOScheme({}),
  identificationType: fieldScheme.String({}),
  identificationValue: fieldScheme.String({}),
})

export default {
  SignUpScheme,
  VerifyEmailScheme,
  ChangePasswdScheme,
  LoginScheme,
  ForgetPasswdScheme,
  ForgetPasswdUpdateScheme,
  OnboardingScheme
}