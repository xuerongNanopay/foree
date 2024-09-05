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
  LastName: fieldScheme.NameScheme({}),
  address1: string().trim().required(),
  address2: string(),
  city: string().trim().required(),
  province: string().trim().required(),
  country: string().trim().required(),
  phoneNumber: string().trim().required(),
  dob: fieldScheme.DateOnlyScheme({}),
  nationality: string().trim().required(),
  identificationType: string().trim().required(),
  identificationValue: string().trim().required(),
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