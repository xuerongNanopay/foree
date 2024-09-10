import axios from 'axios'

class AuthService {
  constructor() {
  }

  async login(req) {
    // this.#axiosClient.post('')
    // this.#axiosClient.post("/")
  }

  signIn() {

  }

  signUp() {

  }


  async forgetPassword(req, {signal}={signal}) {
    return await axios.post("/forget_password", req, {signal})
  }

  onboard() {

  }
}

export default AuthService