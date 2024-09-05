import Joi from 'joi'

import fieldScheme from '../constants/validation_sheme'

const SignUpScheme = Joi.object({
  email: fieldScheme.EmailScheme(),
  password: fieldScheme.PasswordScheme(),
})

const VerifyEmailScheme = Joi.object({
  code: fieldScheme.NumberOnlyScheme({min:6, max:6}),
})

const ChangePasswdScheme = Joi.object({
  oldPassword: Joi.string().required(),
  password: fieldScheme.PasswordScheme(),
})

const LoginScheme = Joi.object({
  email: fieldScheme.EmailScheme(),
  password: fieldScheme.PasswordScheme(),
})

const ForgetPasswdScheme = Joi.object({
  email: fieldScheme.EmailScheme(),
})

const ForgetPasswdUpdateScheme = Joi.object({
  retrieveCode: Joi.string().required(),
  newPassword: fieldScheme.PasswordScheme(),
})


export default {
  SignUpScheme,
  VerifyEmailScheme,
  ChangePasswdScheme,
  LoginScheme,
  ForgetPasswdScheme,
  ForgetPasswdUpdateScheme
}