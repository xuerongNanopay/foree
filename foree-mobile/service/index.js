import AuthService from "./auth_services"
import {
  SignUpScheme,
  VerifyEmailScheme,
  ChangePasswdScheme,
  LoginScheme,
  ForgetPasswdScheme,
  ForgetPasswdUpdateScheme,
  OnboardingScheme
} from "./auth_payload"

const authService = new AuthService()

export {
  authService,
  SignUpScheme,
  VerifyEmailScheme,
  ChangePasswdScheme,
  LoginScheme,
  ForgetPasswdScheme,
  ForgetPasswdUpdateScheme,
  OnboardingScheme
}