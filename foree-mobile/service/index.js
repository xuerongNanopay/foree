import axios from 'axios'

import AuthService from "./auth_services"
import authPayload from "./auth_payload"

// Config axios
axios.defaults.baseURL = 'http://192.168.2.30:8080/app/v1'
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    //Need text
    // return Promise.resolve(error)
    return Promise.reject(error)
  }
)

const authService = new AuthService()

export {
  authService,
  authPayload
}