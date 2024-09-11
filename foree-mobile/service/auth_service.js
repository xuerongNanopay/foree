import axios from 'axios'

class AuthService {
  constructor() {
  }

  async login(req, {signal}={signal}) {
    return await axios.post("/login", req, {signal})
  }

  async signUp(req, {signal}={signal}) {
    return await axios.post("/sign_up", req, {signal})
  }

  async verifyEmail(req, {signal}={signal}) {
    return await axios.post("/verify_email", req, {signal})
  }

  async resendCode({signal}={signal}) {
    return await axios.get("/resend_code", {signal})
  }

  async forgetPassword(req, {signal}={signal}) {
    return await axios.post("/forget_password", req, {signal})
  }

  async forgetPasswordVerify(req, {signal}={signal}) {
    return await axios.post("/forget_password_verify", req, {signal})
  }

  async forgetPasswordUpdate(req, {signal}={signal}) {
    return await axios.post("/forget_password_update", req, {signal})
  }

  async onboard(req, {signal}={signal}) {
    return await axios.post("/onboard", req, {signal})
  }

  async resendCode({signal}={signal}) {
    return await axios.get("/user", {signal})
  }

}

export default AuthService