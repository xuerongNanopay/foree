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
    try {
      const resp = await axios.post("/forget_password", req, {signal})
      const data = resp.data
      return data
    } catch (err) {
      console.log('catch', err.response.status)
      throw err
    }
  }

  onboard() {

  }
}

export default AuthService