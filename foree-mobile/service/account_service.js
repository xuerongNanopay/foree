import axios from 'axios'

class AccountService {
  async getRate(req, {signal}={signal}) {
    return await axios.post("/rate", req, {signal})
  }
  async freeRate(req, {signal}={signal}) {
    return await axios.post("/free_quote", req, {signal})
  }
}