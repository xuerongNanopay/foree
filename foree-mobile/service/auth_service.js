import axios from 'axios'

class AuthService {
  constructor() {
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
    return await axios.post("/sign_up", req, {signal})
  }

  onboard() {

  }
}

export default AuthService