import axios from 'axios'

class AccountService {
  async verifyContact(req, {signal}={signal}) {
    return await axios.post("/verify_contact_account", req, {signal})
  }
}