import axios from 'axios'

class AuthService {
  constructor() {
  }

  async signUp(req, {signal}={signal}) {
    return await axios.post("/sign_up", req, {signal})
  }

  async forgetPassword(req, {signal}={signal}) {
    return await axios.post("/sign_up", req, {signal})
  }

  onboard() {

  }
}

export default AuthService