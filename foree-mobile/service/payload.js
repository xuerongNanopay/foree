import Joi from 'joi'

import fieldScheme from '../constants/validation_sheme'

const SignUpScheme = Joi.object({
  email: fieldScheme.EmailScheme({}),
  password: fieldScheme.PasswordScheme({}),
})

const VerifyEmailScheme = Joi.object({
  code: fieldScheme.NumberOnlyScheme({min:6, max:6}),
})

const ChangePasswdScheme = Joi.object({
  oldPassword: Joi.string().required(),
  password: fieldScheme.PasswordScheme({}),
})

const LoginScheme = Joi.object({
  email: fieldScheme.EmailScheme({}),
  password: fieldScheme.PasswordScheme({}),
})

const ForgetPasswdScheme = Joi.object({
  email: fieldScheme.EmailScheme({}),
})

const ForgetPasswdUpdateScheme = Joi.object({
  retrieveCode: Joi.string().required(),
  newPassword: fieldScheme.PasswordScheme({}),
})

const OnboardingScheme = Joi.object({
  firstName: fieldScheme.NameScheme({}),
  middleName: fieldScheme.NameScheme({required: false}),
  LastName: fieldScheme.NameScheme({}),
  address1: Joi.string().required(),
  address2: Joi.string(),
  city: Joi.string().required(),
  province: Joi.string().required().custom((value, helper) => {
    console.log("Onboarding Scheme", helper.original)
    return true
  }),
  country: Joi.string().required(),
  phoneNumber: Joi.string().required(),
  dob: fieldScheme.DateOnlyScheme({}),
  nationality: Joi.string().required(),
  identificationType: Joi.string().required(),
  identificationValue: Joi.string().required(),
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