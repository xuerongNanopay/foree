import axios from 'axios'

class AuthService {
  #serviceConfig
  #axiosClient


  constructor() {
    this.#serviceConfig = {
      endPoint: "http://localhost:8080/"
    }
    this.#axiosClient = axios.create({
      baseURL: 'http://localhost:8080/app/v1'
    })

    this.#axiosClient.interceptors.response.use(
      (response) => response,
      (error) => {
        //Need text
        console.log(error.response)
        return Promise.reject(error)
      }
    )
  }

  async login(req) {
    // this.#axiosClient.post('')
    // this.#axiosClient.post("/")
  }

  signIn() {

  }

  signUp() {

  }


  async forgetPassword(req, {signal=new axios.AbortController()}) {
    try {
      const resp = await this.#axiosClient.post("/forget_password", req, {signal})
      const data = resp.data
      return data
    } catch (error) {
      return error.response
    }
  }

  onboard() {

  }
}

export default AuthService